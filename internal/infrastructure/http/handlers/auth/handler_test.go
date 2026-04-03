package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"backend-typing-trainer/internal/domain/models"
	infraLogger "backend-typing-trainer/internal/infrastructure/logger"
	"backend-typing-trainer/internal/utils"
	"backend-typing-trainer/mocks"
)

func newTestAuthHandler(t *testing.T, svc *mocks.AuthInputPort) *Handler {
	t.Helper()

	return NewHandler(infraLogger.New("dev"), svc)
}

func decodeErrorResponse(t *testing.T, rr *httptest.ResponseRecorder) utils.ErrorResponse {
	t.Helper()

	var resp utils.ErrorResponse
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &resp))
	return resp
}

func TestHandler_Register(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		authSvc := mocks.NewAuthInputPort(t)
		h := newTestAuthHandler(t, authSvc)

		userID := uuid.MustParse("22222222-2222-2222-2222-222222222222")
		createdAt := time.Date(2026, time.March, 22, 21, 0, 0, 0, time.UTC)

		authSvc.EXPECT().Register(mock.Anything, "player_one", "secret", models.UserRoleUser).Return(&models.User{
			ID:        userID,
			Login:     "player_one",
			Role:      models.UserRoleUser,
			CreatedAt: createdAt,
		}, nil).Once()

		r := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(`{"login":"player_one","password":"secret"}`))
		rr := httptest.NewRecorder()

		h.Register(rr, r)

		require.Equal(t, http.StatusCreated, rr.Code)
		var resp registerResponse
		require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &resp))
		require.Equal(t, userID, resp.User.ID)
		require.Equal(t, "player_one", resp.User.Login)
		require.Equal(t, models.UserRoleUser, resp.User.Role)
		require.NotNil(t, resp.User.CreatedAt)
		require.True(t, resp.User.CreatedAt.Equal(createdAt))
	})

	t.Run("duplicate login returns 409 login exists", func(t *testing.T) {
		authSvc := mocks.NewAuthInputPort(t)
		h := newTestAuthHandler(t, authSvc)

		authSvc.EXPECT().Register(mock.Anything, "player_one", "secret", models.UserRoleUser).
			Return(nil, utils.ErrUserLoginExists).Once()

		r := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(`{"login":"player_one","password":"secret"}`))
		rr := httptest.NewRecorder()

		h.Register(rr, r)

		require.Equal(t, http.StatusConflict, rr.Code)
		errResp := decodeErrorResponse(t, rr)
		require.Equal(t, utils.ErrorCodeLoginExists, errResp.Error.Code)
		require.Equal(t, "login already exists", errResp.Error.Message)
	})

	t.Run("unknown field returns 400", func(t *testing.T) {
		authSvc := mocks.NewAuthInputPort(t)
		h := newTestAuthHandler(t, authSvc)

		r := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(`{"login":"player_one","password":"secret","extra":"x"}`))
		rr := httptest.NewRecorder()

		h.Register(rr, r)

		require.Equal(t, http.StatusBadRequest, rr.Code)
		errResp := decodeErrorResponse(t, rr)
		require.Equal(t, utils.ErrorCodeInvalidRequest, errResp.Error.Code)
		require.Equal(t, "invalid request", errResp.Error.Message)
	})

	t.Run("nil service returns 500", func(t *testing.T) {
		h := NewHandler(infraLogger.New("dev"), nil)

		r := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(`{"login":"player_one","password":"secret"}`))
		rr := httptest.NewRecorder()

		h.Register(rr, r)

		require.Equal(t, http.StatusInternalServerError, rr.Code)
		errResp := decodeErrorResponse(t, rr)
		require.Equal(t, utils.ErrorCodeInternalError, errResp.Error.Code)
		require.Equal(t, "internal server error", errResp.Error.Message)
	})
}

func TestHandler_Login(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		authSvc := mocks.NewAuthInputPort(t)
		h := newTestAuthHandler(t, authSvc)

		authSvc.EXPECT().Login(mock.Anything, "player_one", "secret").Return("jwt-token", nil).Once()

		r := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(`{"login":"player_one","password":"secret"}`))
		rr := httptest.NewRecorder()

		h.Login(rr, r)

		require.Equal(t, http.StatusOK, rr.Code)
		var resp LoginResponse
		require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &resp))
		require.Equal(t, "jwt-token", resp.Token)
	})

	t.Run("invalid body returns 401", func(t *testing.T) {
		authSvc := mocks.NewAuthInputPort(t)
		h := newTestAuthHandler(t, authSvc)

		r := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(`{"login":"","password":"secret"}`))
		rr := httptest.NewRecorder()

		h.Login(rr, r)

		require.Equal(t, http.StatusUnauthorized, rr.Code)
		errResp := decodeErrorResponse(t, rr)
		require.Equal(t, utils.ErrorCodeUnauthorized, errResp.Error.Code)
		require.Equal(t, "unauthorized", errResp.Error.Message)
	})

	t.Run("service unauthorized returns 401", func(t *testing.T) {
		authSvc := mocks.NewAuthInputPort(t)
		h := newTestAuthHandler(t, authSvc)

		authSvc.EXPECT().Login(mock.Anything, "player_one", "secret").Return("", utils.ErrUnauthorized).Once()

		r := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(`{"login":"player_one","password":"secret"}`))
		rr := httptest.NewRecorder()

		h.Login(rr, r)

		require.Equal(t, http.StatusUnauthorized, rr.Code)
		errResp := decodeErrorResponse(t, rr)
		require.Equal(t, utils.ErrorCodeUnauthorized, errResp.Error.Code)
		require.Equal(t, "unauthorized", errResp.Error.Message)
	})

	t.Run("service internal error returns 500", func(t *testing.T) {
		authSvc := mocks.NewAuthInputPort(t)
		h := newTestAuthHandler(t, authSvc)

		authSvc.EXPECT().Login(mock.Anything, "player_one", "secret").Return("", errors.New("boom")).Once()

		r := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(`{"login":"player_one","password":"secret"}`))
		rr := httptest.NewRecorder()

		h.Login(rr, r)

		require.Equal(t, http.StatusInternalServerError, rr.Code)
		errResp := decodeErrorResponse(t, rr)
		require.Equal(t, utils.ErrorCodeInternalError, errResp.Error.Code)
		require.Equal(t, "internal server error", errResp.Error.Message)
	})

	t.Run("nil service returns 500", func(t *testing.T) {
		h := NewHandler(infraLogger.New("dev"), nil)

		r := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(`{"login":"player_one","password":"secret"}`))
		rr := httptest.NewRecorder()

		h.Login(rr, r)

		require.Equal(t, http.StatusInternalServerError, rr.Code)
		errResp := decodeErrorResponse(t, rr)
		require.Equal(t, utils.ErrorCodeInternalError, errResp.Error.Code)
		require.Equal(t, "internal server error", errResp.Error.Message)
	})
}
