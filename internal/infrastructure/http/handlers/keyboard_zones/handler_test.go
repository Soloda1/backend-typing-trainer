package keyboardzones

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

func newTestKeyboardZonesHandler(t *testing.T, svc *mocks.KeyboardZonesInputPort) *Handler {
	t.Helper()
	return NewHandler(infraLogger.New("dev"), svc)
}

func decodeErrorResponse(t *testing.T, rr *httptest.ResponseRecorder) utils.ErrorResponse {
	t.Helper()
	var resp utils.ErrorResponse
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &resp))
	return resp
}

func withID(r *http.Request, id string) *http.Request {
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", id)
	ctx := context.WithValue(r.Context(), chi.RouteCtxKey, rctx)
	return r.WithContext(ctx)
}

func withName(r *http.Request, name string) *http.Request {
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("name", name)
	ctx := context.WithValue(r.Context(), chi.RouteCtxKey, rctx)
	return r.WithContext(ctx)
}

func TestHandler_List(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := mocks.NewKeyboardZonesInputPort(t)
		h := newTestKeyboardZonesHandler(t, svc)

		svc.EXPECT().List(mock.Anything, 20, 0).Return([]*models.KeyboardZone{
			{ID: uuid.MustParse("11111111-1111-1111-1111-111111111111"), Name: "en_red", Symbols: "q,a,z"},
		}, nil).Once()

		r := httptest.NewRequest(http.MethodGet, "/keyboard-zones", nil)
		rr := httptest.NewRecorder()

		h.List(rr, r)

		require.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("invalid pagination returns 400", func(t *testing.T) {
		svc := mocks.NewKeyboardZonesInputPort(t)
		h := newTestKeyboardZonesHandler(t, svc)

		r := httptest.NewRequest(http.MethodGet, "/keyboard-zones?limit=bad", nil)
		rr := httptest.NewRecorder()

		h.List(rr, r)

		require.Equal(t, http.StatusBadRequest, rr.Code)
		errResp := decodeErrorResponse(t, rr)
		require.Equal(t, utils.ErrorCodeInvalidRequest, errResp.Error.Code)
	})
}

func TestHandler_GetByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := mocks.NewKeyboardZonesInputPort(t)
		h := newTestKeyboardZonesHandler(t, svc)

		id := uuid.MustParse("22222222-2222-2222-2222-222222222222")
		svc.EXPECT().GetByID(mock.Anything, id).Return(&models.KeyboardZone{ID: id, Name: "en_blue", Symbols: "y,u"}, nil).Once()

		r := httptest.NewRequest(http.MethodGet, "/keyboard-zones/"+id.String(), nil)
		r = withID(r, id.String())
		rr := httptest.NewRecorder()

		h.GetByID(rr, r)

		require.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("invalid id returns 400", func(t *testing.T) {
		svc := mocks.NewKeyboardZonesInputPort(t)
		h := newTestKeyboardZonesHandler(t, svc)

		r := httptest.NewRequest(http.MethodGet, "/keyboard-zones/bad", nil)
		r = withID(r, "bad")
		rr := httptest.NewRecorder()

		h.GetByID(rr, r)

		require.Equal(t, http.StatusBadRequest, rr.Code)
		errResp := decodeErrorResponse(t, rr)
		require.Equal(t, utils.ErrorCodeInvalidRequest, errResp.Error.Code)
	})
}

func TestHandler_GetByName(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := mocks.NewKeyboardZonesInputPort(t)
		h := newTestKeyboardZonesHandler(t, svc)

		svc.EXPECT().GetByName(mock.Anything, "en_red").Return(&models.KeyboardZone{Name: "en_red", Symbols: "q,a,z"}, nil).Once()

		r := httptest.NewRequest(http.MethodGet, "/keyboard-zones/by-name/en_red", nil)
		r = withName(r, "en_red")
		rr := httptest.NewRecorder()

		h.GetByName(rr, r)

		require.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("service internal error returns 500", func(t *testing.T) {
		svc := mocks.NewKeyboardZonesInputPort(t)
		h := newTestKeyboardZonesHandler(t, svc)

		svc.EXPECT().GetByName(mock.Anything, "en_red").Return(nil, errors.New("boom")).Once()

		r := httptest.NewRequest(http.MethodGet, "/keyboard-zones/by-name/en_red", nil)
		r = withName(r, "en_red")
		rr := httptest.NewRecorder()

		h.GetByName(rr, r)

		require.Equal(t, http.StatusInternalServerError, rr.Code)
		errResp := decodeErrorResponse(t, rr)
		require.Equal(t, utils.ErrorCodeInternalError, errResp.Error.Code)
	})
}
