package middlewares

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"backend-typing-trainer/internal/domain/models"
	input "backend-typing-trainer/internal/domain/ports/input"
	infraLogger "backend-typing-trainer/internal/infrastructure/logger"
	"backend-typing-trainer/internal/utils"
)

func TestRequireRoles(t *testing.T) {
	t.Run("missing actor returns 401", func(t *testing.T) {
		mw := RequireRoles(infraLogger.New("dev"), models.UserRoleAdmin)
		h := mw(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
			t.Fatal("next handler should not be called")
		}))

		req := httptest.NewRequest(http.MethodGet, "/admin", http.NoBody)
		rr := httptest.NewRecorder()

		h.ServeHTTP(rr, req)

		require.Equal(t, http.StatusUnauthorized, rr.Code)
		var errResp utils.ErrorResponse
		require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &errResp))
		require.Equal(t, utils.ErrorCodeUnauthorized, errResp.Error.Code)
	})

	t.Run("forbidden role returns 403", func(t *testing.T) {
		mw := RequireRoles(infraLogger.New("dev"), models.UserRoleAdmin)
		h := mw(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
			t.Fatal("next handler should not be called")
		}))

		actor := input.Actor{UserID: uuid.MustParse("22222222-2222-2222-2222-222222222222"), Role: models.UserRoleUser}
		req := httptest.NewRequest(http.MethodGet, "/admin", http.NoBody)
		req = req.WithContext(utils.WithActor(req.Context(), actor))
		rr := httptest.NewRecorder()

		h.ServeHTTP(rr, req)

		require.Equal(t, http.StatusForbidden, rr.Code)
		var errResp utils.ErrorResponse
		require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &errResp))
		require.Equal(t, utils.ErrorCodeForbidden, errResp.Error.Code)
	})

	t.Run("allowed role passes next", func(t *testing.T) {
		mw := RequireRoles(infraLogger.New("dev"), models.UserRoleAdmin, models.UserRoleUser)
		called := false
		h := mw(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			called = true
			w.WriteHeader(http.StatusNoContent)
		}))

		actor := input.Actor{UserID: uuid.MustParse("22222222-2222-2222-2222-222222222222"), Role: models.UserRoleUser}
		req := httptest.NewRequest(http.MethodGet, "/user", http.NoBody)
		req = req.WithContext(utils.WithActor(req.Context(), actor))
		rr := httptest.NewRecorder()

		h.ServeHTTP(rr, req)

		require.True(t, called)
		require.Equal(t, http.StatusNoContent, rr.Code)
	})
}

func TestRequireSelfOrRoles(t *testing.T) {
	log := infraLogger.New("dev")

	newRouter := func(actor *input.Actor) *chi.Mux {
		r := chi.NewRouter()
		r.With(RequireSelfOrRoles(log, "user_id", models.UserRoleAdmin)).Get("/statistics/users/{user_id}", func(w http.ResponseWriter, req *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		})

		if actor == nil {
			return r
		}

		wrapped := chi.NewRouter()
		wrapped.Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				next.ServeHTTP(w, req.WithContext(utils.WithActor(req.Context(), *actor)))
			})
		})
		wrapped.Mount("/", r)
		return wrapped
	}

	t.Run("self access allowed", func(t *testing.T) {
		userID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
		r := newRouter(&input.Actor{UserID: userID, Role: models.UserRoleUser})

		req := httptest.NewRequest(http.MethodGet, "/statistics/users/"+userID.String(), http.NoBody)
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		require.Equal(t, http.StatusNoContent, rr.Code)
	})

	t.Run("admin access allowed", func(t *testing.T) {
		targetID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
		r := newRouter(&input.Actor{UserID: uuid.MustParse("22222222-2222-2222-2222-222222222222"), Role: models.UserRoleAdmin})

		req := httptest.NewRequest(http.MethodGet, "/statistics/users/"+targetID.String(), http.NoBody)
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		require.Equal(t, http.StatusNoContent, rr.Code)
	})

	t.Run("foreign user forbidden", func(t *testing.T) {
		targetID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
		r := newRouter(&input.Actor{UserID: uuid.MustParse("33333333-3333-3333-3333-333333333333"), Role: models.UserRoleUser})

		req := httptest.NewRequest(http.MethodGet, "/statistics/users/"+targetID.String(), http.NoBody)
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		require.Equal(t, http.StatusForbidden, rr.Code)
		var errResp utils.ErrorResponse
		require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &errResp))
		require.Equal(t, utils.ErrorCodeForbidden, errResp.Error.Code)
	})

	t.Run("invalid user id returns 400", func(t *testing.T) {
		r := newRouter(&input.Actor{UserID: uuid.MustParse("33333333-3333-3333-3333-333333333333"), Role: models.UserRoleAdmin})

		req := httptest.NewRequest(http.MethodGet, "/statistics/users/not-uuid", http.NoBody)
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		require.Equal(t, http.StatusBadRequest, rr.Code)
		var errResp utils.ErrorResponse
		require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &errResp))
		require.Equal(t, utils.ErrorCodeInvalidRequest, errResp.Error.Code)
	})

	t.Run("missing actor returns 401", func(t *testing.T) {
		userID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
		r := newRouter(nil)

		req := httptest.NewRequest(http.MethodGet, "/statistics/users/"+userID.String(), http.NoBody)
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		require.Equal(t, http.StatusUnauthorized, rr.Code)
		var errResp utils.ErrorResponse
		require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &errResp))
		require.Equal(t, utils.ErrorCodeUnauthorized, errResp.Error.Code)
	})
}
