package middlewares

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"backend-typing-trainer/internal/domain/models"
	jwtport "backend-typing-trainer/internal/domain/ports/output/jwt"
	infraLogger "backend-typing-trainer/internal/infrastructure/logger"
	"backend-typing-trainer/internal/utils"
	"backend-typing-trainer/mocks"
)

func TestExtractBearerToken(t *testing.T) {
	tests := []struct {
		name   string
		header string
		token  string
		ok     bool
	}{
		{name: "empty header", header: "", token: "", ok: false},
		{name: "without space", header: "BearerToken", token: "", ok: false},
		{name: "wrong scheme", header: "Basic abc", token: "", ok: false},
		{name: "empty token", header: "Bearer ", token: "", ok: false},
		{name: "valid token", header: "Bearer abc.def.ghi", token: "abc.def.ghi", ok: true},
		{name: "case-insensitive scheme", header: "bearer token-1", token: "token-1", ok: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotToken, gotOK := extractBearerToken(tt.header)
			require.Equal(t, tt.token, gotToken)
			require.Equal(t, tt.ok, gotOK)
		})
	}
}

func TestAuthMiddleware(t *testing.T) {
	t.Run("nil token manager returns 500", func(t *testing.T) {
		mw := AuthMiddleware(nil, infraLogger.New("dev"))
		h := mw(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
			t.Fatal("next handler should not be called")
		}))

		req := httptest.NewRequest(http.MethodGet, "/protected", http.NoBody)
		rr := httptest.NewRecorder()

		h.ServeHTTP(rr, req)

		require.Equal(t, http.StatusInternalServerError, rr.Code)
		var errResp utils.ErrorResponse
		require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &errResp))
		require.Equal(t, utils.ErrorCodeInternalError, errResp.Error.Code)
	})

	t.Run("missing auth header returns 401", func(t *testing.T) {
		tm := mocks.NewTokenManager(t)
		mw := AuthMiddleware(tm, infraLogger.New("dev"))
		h := mw(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
			t.Fatal("next handler should not be called")
		}))

		req := httptest.NewRequest(http.MethodGet, "/protected", http.NoBody)
		rr := httptest.NewRecorder()

		h.ServeHTTP(rr, req)

		require.Equal(t, http.StatusUnauthorized, rr.Code)
		var errResp utils.ErrorResponse
		require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &errResp))
		require.Equal(t, utils.ErrorCodeUnauthorized, errResp.Error.Code)
	})

	t.Run("parse token error maps to 401", func(t *testing.T) {
		tm := mocks.NewTokenManager(t)
		tm.EXPECT().ParseToken("bad-token").Return(nil, utils.ErrUnauthorized).Once()

		mw := AuthMiddleware(tm, infraLogger.New("dev"))
		h := mw(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
			t.Fatal("next handler should not be called")
		}))

		req := httptest.NewRequest(http.MethodGet, "/protected", http.NoBody)
		req.Header.Set("Authorization", "Bearer bad-token")
		rr := httptest.NewRecorder()

		h.ServeHTTP(rr, req)

		require.Equal(t, http.StatusUnauthorized, rr.Code)
		var errResp utils.ErrorResponse
		require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &errResp))
		require.Equal(t, utils.ErrorCodeUnauthorized, errResp.Error.Code)
	})

	t.Run("unexpected parse token error maps to 500", func(t *testing.T) {
		tm := mocks.NewTokenManager(t)
		tm.EXPECT().ParseToken("bad-token").Return(nil, errors.New("boom")).Once()

		mw := AuthMiddleware(tm, infraLogger.New("dev"))
		h := mw(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
			t.Fatal("next handler should not be called")
		}))

		req := httptest.NewRequest(http.MethodGet, "/protected", http.NoBody)
		req.Header.Set("Authorization", "Bearer bad-token")
		rr := httptest.NewRecorder()

		h.ServeHTTP(rr, req)

		require.Equal(t, http.StatusInternalServerError, rr.Code)
		var errResp utils.ErrorResponse
		require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &errResp))
		require.Equal(t, utils.ErrorCodeInternalError, errResp.Error.Code)
	})

	t.Run("success passes actor to next", func(t *testing.T) {
		tm := mocks.NewTokenManager(t)
		userID := uuid.MustParse("22222222-2222-2222-2222-222222222222")
		tm.EXPECT().ParseToken("good-token").Return(&jwtport.TokenClaims{UserID: userID, Role: models.UserRoleUser}, nil).Once()

		called := false
		mw := AuthMiddleware(tm, infraLogger.New("dev"))
		h := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			actor, ok := utils.ActorFromContext(r.Context())
			require.True(t, ok)
			require.Equal(t, userID, actor.UserID)
			require.Equal(t, models.UserRoleUser, actor.Role)
			w.WriteHeader(http.StatusNoContent)
		}))

		req := httptest.NewRequest(http.MethodGet, "/protected", http.NoBody)
		req.Header.Set("Authorization", "Bearer good-token")
		rr := httptest.NewRecorder()

		h.ServeHTTP(rr, req)

		require.True(t, called)
		require.Equal(t, http.StatusNoContent, rr.Code)
	})
}
