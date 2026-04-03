package middlewares

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

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
