package exercises

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
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

func newTestExercisesHandler(t *testing.T, svc *mocks.ExercisesInputPort) *Handler {
	t.Helper()
	return NewHandler(infraLogger.New("dev"), svc)
}

func decodeErr(t *testing.T, rr *httptest.ResponseRecorder) utils.ErrorResponse {
	t.Helper()
	var resp utils.ErrorResponse
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &resp))
	return resp
}

func withRouteID(r *http.Request, id string) *http.Request {
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", id)
	ctx := context.WithValue(r.Context(), chi.RouteCtxKey, rctx)
	return r.WithContext(ctx)
}

func TestHandler_Create(t *testing.T) {
	svc := mocks.NewExercisesInputPort(t)
	h := newTestExercisesHandler(t, svc)

	createdID := uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa")
	svc.EXPECT().Create(mock.Anything, mock.AnythingOfType("*models.Exercise")).
		Run(func(_ context.Context, ex *models.Exercise) { ex.ID = createdID }).
		Return(nil).Once()

	r := httptest.NewRequest(http.MethodPost, "/exercises", strings.NewReader(`{"text":"hello","level_id":"11111111-1111-1111-1111-111111111111"}`))
	rr := httptest.NewRecorder()

	h.Create(rr, r)

	require.Equal(t, http.StatusCreated, rr.Code)

	r = httptest.NewRequest(http.MethodPost, "/exercises", strings.NewReader(`{"text":"hello","level_id":"11111111-1111-1111-1111-111111111111","extra":true}`))
	rr = httptest.NewRecorder()
	h.Create(rr, r)
	require.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestHandler_GetByID(t *testing.T) {
	svc := mocks.NewExercisesInputPort(t)
	h := newTestExercisesHandler(t, svc)

	id := uuid.MustParse("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb")
	svc.EXPECT().GetByID(mock.Anything, id).Return(&models.Exercise{ID: id}, nil).Once()

	r := httptest.NewRequest(http.MethodGet, "/exercises/"+id.String(), nil)
	r = withRouteID(r, id.String())
	rr := httptest.NewRecorder()
	h.GetByID(rr, r)
	require.Equal(t, http.StatusOK, rr.Code)

	r = httptest.NewRequest(http.MethodGet, "/exercises/bad", nil)
	r = withRouteID(r, "bad")
	rr = httptest.NewRecorder()
	h.GetByID(rr, r)
	require.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestHandler_Delete(t *testing.T) {
	svc := mocks.NewExercisesInputPort(t)
	h := newTestExercisesHandler(t, svc)

	id := uuid.MustParse("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb")

	svc.EXPECT().Delete(mock.Anything, id).Return(errors.New("boom")).Once()
	r := httptest.NewRequest(http.MethodDelete, "/exercises/"+id.String(), nil)
	r = withRouteID(r, id.String())
	rr := httptest.NewRecorder()
	h.Delete(rr, r)
	require.Equal(t, http.StatusInternalServerError, rr.Code)
	errResp := decodeErr(t, rr)
	require.Equal(t, utils.ErrorCodeInternalError, errResp.Error.Code)

	svc.EXPECT().Delete(mock.Anything, id).Return(nil).Once()
	r = httptest.NewRequest(http.MethodDelete, "/exercises/"+id.String(), nil)
	r = withRouteID(r, id.String())
	rr = httptest.NewRecorder()
	h.Delete(rr, r)
	require.Equal(t, http.StatusNoContent, rr.Code)
}

func TestHandler_List(t *testing.T) {
	svc := mocks.NewExercisesInputPort(t)
	h := newTestExercisesHandler(t, svc)

	r := httptest.NewRequest(http.MethodGet, "/exercises?limit=bad", nil)
	rr := httptest.NewRecorder()
	h.List(rr, r)
	require.Equal(t, http.StatusBadRequest, rr.Code)

	svc.EXPECT().List(mock.Anything, 20, 0).Return([]*models.Exercise{{ID: uuid.MustParse("dddddddd-dddd-dddd-dddd-dddddddddddd")}}, nil).Once()
	r = httptest.NewRequest(http.MethodGet, "/exercises", nil)
	rr = httptest.NewRecorder()
	h.List(rr, r)
	require.Equal(t, http.StatusOK, rr.Code)
}

func TestHandler_Update(t *testing.T) {
	svc := mocks.NewExercisesInputPort(t)
	h := newTestExercisesHandler(t, svc)

	id := uuid.MustParse("cccccccc-cccc-cccc-cccc-cccccccccccc")
	r := httptest.NewRequest(http.MethodPatch, "/exercises/"+id.String(), strings.NewReader(`{"text":"x","level_id":"11111111-1111-1111-1111-111111111111","extra":true}`))
	r = withRouteID(r, id.String())
	rr := httptest.NewRecorder()
	h.Update(rr, r)
	require.Equal(t, http.StatusBadRequest, rr.Code)

	svc.EXPECT().Update(mock.Anything, mock.AnythingOfType("*models.Exercise")).Return(nil).Once()
	r = httptest.NewRequest(http.MethodPatch, "/exercises/"+id.String(), strings.NewReader(`{"text":"updated","level_id":"11111111-1111-1111-1111-111111111111"}`))
	r = withRouteID(r, id.String())
	rr = httptest.NewRecorder()
	h.Update(rr, r)
	require.Equal(t, http.StatusOK, rr.Code)
}
