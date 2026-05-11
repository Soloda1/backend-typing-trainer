package users

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"backend-typing-trainer/internal/domain/models"
	infraLogger "backend-typing-trainer/internal/infrastructure/logger"
	"backend-typing-trainer/internal/utils"
	"backend-typing-trainer/mocks"
)

func newTestUsersHandler(t *testing.T, svc *mocks.UsersInputPort) *Handler {
	t.Helper()
	return NewHandler(svc, infraLogger.New("dev"))
}

func decodeUsersErrorResponse(t *testing.T, rr *httptest.ResponseRecorder) utils.ErrorResponse {
	t.Helper()
	var resp utils.ErrorResponse
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &resp))
	return resp
}

func withRouteParam(r *http.Request, key, value string) *http.Request {
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add(key, value)
	ctx := context.WithValue(r.Context(), chi.RouteCtxKey, rctx)
	return r.WithContext(ctx)
}

func TestHandler_GetByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := mocks.NewUsersInputPort(t)
		h := newTestUsersHandler(t, svc)

		id := uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa")
		svc.EXPECT().GetByID(mock.Anything, id).Return(&models.User{ID: id}, nil).Once()

		r := httptest.NewRequest(http.MethodGet, "/users/"+id.String(), nil)
		r = withRouteParam(r, "id", id.String())
		rr := httptest.NewRecorder()

		h.GetByID(rr, r)

		require.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("invalid id returns 400", func(t *testing.T) {
		svc := mocks.NewUsersInputPort(t)
		h := newTestUsersHandler(t, svc)

		r := httptest.NewRequest(http.MethodGet, "/users/bad", nil)
		r = withRouteParam(r, "id", "bad")
		rr := httptest.NewRecorder()

		h.GetByID(rr, r)

		require.Equal(t, http.StatusBadRequest, rr.Code)
		errResp := decodeUsersErrorResponse(t, rr)
		require.Equal(t, utils.ErrorCodeInvalidRequest, errResp.Error.Code)
	})
}

func TestHandler_GetByLogin(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := mocks.NewUsersInputPort(t)
		h := newTestUsersHandler(t, svc)

		svc.EXPECT().GetByLogin(mock.Anything, "admin").Return(&models.User{Login: "admin"}, nil).Once()

		r := httptest.NewRequest(http.MethodGet, "/users/login/admin", nil)
		r = withRouteParam(r, "login", "admin")
		rr := httptest.NewRecorder()

		h.GetByLogin(rr, r)

		require.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("empty login returns 400", func(t *testing.T) {
		svc := mocks.NewUsersInputPort(t)
		h := newTestUsersHandler(t, svc)

		r := httptest.NewRequest(http.MethodGet, "/users/login/", nil)
		r = withRouteParam(r, "login", " ")
		rr := httptest.NewRecorder()

		h.GetByLogin(rr, r)

		require.Equal(t, http.StatusBadRequest, rr.Code)
		errResp := decodeUsersErrorResponse(t, rr)
		require.Equal(t, utils.ErrorCodeInvalidRequest, errResp.Error.Code)
	})

	t.Run("service error returns 500", func(t *testing.T) {
		svc := mocks.NewUsersInputPort(t)
		h := newTestUsersHandler(t, svc)

		svc.EXPECT().GetByLogin(mock.Anything, "admin").Return(nil, errors.New("boom")).Once()

		r := httptest.NewRequest(http.MethodGet, "/users/login/admin", nil)
		r = withRouteParam(r, "login", "admin")
		rr := httptest.NewRecorder()

		h.GetByLogin(rr, r)

		require.Equal(t, http.StatusInternalServerError, rr.Code)
		errResp := decodeUsersErrorResponse(t, rr)
		require.Equal(t, utils.ErrorCodeInternalError, errResp.Error.Code)
	})
}

func TestHandler_List(t *testing.T) {
	t.Run("success with defaults", func(t *testing.T) {
		svc := mocks.NewUsersInputPort(t)
		h := newTestUsersHandler(t, svc)

		svc.EXPECT().List(mock.Anything, 10, 0).Return([]*models.User{}, nil).Once()

		r := httptest.NewRequest(http.MethodGet, "/users", nil)
		rr := httptest.NewRecorder()

		h.List(rr, r)

		require.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("invalid pagination returns 400", func(t *testing.T) {
		svc := mocks.NewUsersInputPort(t)
		h := newTestUsersHandler(t, svc)

		r := httptest.NewRequest(http.MethodGet, "/users?limit=bad", nil)
		rr := httptest.NewRecorder()

		h.List(rr, r)

		require.Equal(t, http.StatusBadRequest, rr.Code)
		errResp := decodeUsersErrorResponse(t, rr)
		require.Equal(t, utils.ErrorCodeInvalidRequest, errResp.Error.Code)
	})
}
