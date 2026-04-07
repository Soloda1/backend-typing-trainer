package difficultylevels

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

func newTestDifficultyLevelsHandler(t *testing.T, svc *mocks.DifficultyLevelsInputPort) *Handler {
	t.Helper()

	return NewHandler(infraLogger.New("dev"), svc)
}

func decodeDifficultyLevelsErrorResponse(t *testing.T, rr *httptest.ResponseRecorder) utils.ErrorResponse {
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
	t.Run("success", func(t *testing.T) {
		svc := mocks.NewDifficultyLevelsInputPort(t)
		h := newTestDifficultyLevelsHandler(t, svc)

		createdID := uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa")
		svc.EXPECT().Create(mock.Anything, mock.AnythingOfType("*models.DifficultyLevel")).
			Run(func(_ context.Context, level *models.DifficultyLevel) {
				level.ID = createdID
			}).
			Return(nil).Once()

		r := httptest.NewRequest(http.MethodPost, "/difficulty-levels", strings.NewReader(`{"allowed_mistakes":2,"key_press_time":1.5,"min_exercise_length":5,"max_exercise_length":25,"keyboard_zone_ids":["11111111-1111-1111-1111-111111111111"]}`))
		rr := httptest.NewRecorder()

		h.Create(rr, r)

		require.Equal(t, http.StatusCreated, rr.Code)
		var resp difficultyLevelSingleResponse
		require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &resp))
		require.Equal(t, createdID, resp.DifficultyLevel.ID)
	})

	t.Run("unknown field returns 400", func(t *testing.T) {
		svc := mocks.NewDifficultyLevelsInputPort(t)
		h := newTestDifficultyLevelsHandler(t, svc)

		r := httptest.NewRequest(http.MethodPost, "/difficulty-levels", strings.NewReader(`{"allowed_mistakes":2,"key_press_time":1.5,"min_exercise_length":5,"max_exercise_length":25,"keyboard_zone_ids":["11111111-1111-1111-1111-111111111111"],"extra":true}`))
		rr := httptest.NewRecorder()

		h.Create(rr, r)

		require.Equal(t, http.StatusBadRequest, rr.Code)
		errResp := decodeDifficultyLevelsErrorResponse(t, rr)
		require.Equal(t, utils.ErrorCodeInvalidRequest, errResp.Error.Code)
		svc.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
	})
}

func TestHandler_GetByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := mocks.NewDifficultyLevelsInputPort(t)
		h := newTestDifficultyLevelsHandler(t, svc)

		id := uuid.MustParse("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb")
		svc.EXPECT().GetByID(mock.Anything, id).Return(&models.DifficultyLevel{ID: id}, nil).Once()

		r := httptest.NewRequest(http.MethodGet, "/difficulty-levels/"+id.String(), nil)
		r = withRouteID(r, id.String())
		rr := httptest.NewRecorder()

		h.GetByID(rr, r)

		require.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("invalid id returns 400", func(t *testing.T) {
		svc := mocks.NewDifficultyLevelsInputPort(t)
		h := newTestDifficultyLevelsHandler(t, svc)

		r := httptest.NewRequest(http.MethodGet, "/difficulty-levels/bad", nil)
		r = withRouteID(r, "bad")
		rr := httptest.NewRecorder()

		h.GetByID(rr, r)

		require.Equal(t, http.StatusBadRequest, rr.Code)
		errResp := decodeDifficultyLevelsErrorResponse(t, rr)
		require.Equal(t, utils.ErrorCodeInvalidRequest, errResp.Error.Code)
	})
}

func TestHandler_List(t *testing.T) {
	t.Run("success with defaults", func(t *testing.T) {
		svc := mocks.NewDifficultyLevelsInputPort(t)
		h := newTestDifficultyLevelsHandler(t, svc)

		svc.EXPECT().List(mock.Anything, 20, 0).Return([]*models.DifficultyLevel{}, nil).Once()

		r := httptest.NewRequest(http.MethodGet, "/difficulty-levels", nil)
		rr := httptest.NewRecorder()

		h.List(rr, r)

		require.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("invalid query params returns 400", func(t *testing.T) {
		svc := mocks.NewDifficultyLevelsInputPort(t)
		h := newTestDifficultyLevelsHandler(t, svc)

		r := httptest.NewRequest(http.MethodGet, "/difficulty-levels?limit=bad", nil)
		rr := httptest.NewRecorder()

		h.List(rr, r)

		require.Equal(t, http.StatusBadRequest, rr.Code)
		errResp := decodeDifficultyLevelsErrorResponse(t, rr)
		require.Equal(t, utils.ErrorCodeInvalidRequest, errResp.Error.Code)
	})
}

func TestHandler_UpdateAndDelete(t *testing.T) {
	t.Run("update success", func(t *testing.T) {
		svc := mocks.NewDifficultyLevelsInputPort(t)
		h := newTestDifficultyLevelsHandler(t, svc)

		id := uuid.MustParse("cccccccc-cccc-cccc-cccc-cccccccccccc")
		svc.EXPECT().Update(mock.Anything, mock.AnythingOfType("*models.DifficultyLevel")).Return(nil).Once()

		r := httptest.NewRequest(http.MethodPatch, "/difficulty-levels/"+id.String(), strings.NewReader(`{"allowed_mistakes":1,"key_press_time":1.1,"min_exercise_length":4,"max_exercise_length":20,"keyboard_zone_ids":["11111111-1111-1111-1111-111111111111"]}`))
		r = withRouteID(r, id.String())
		rr := httptest.NewRecorder()

		h.Update(rr, r)

		require.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("delete success", func(t *testing.T) {
		svc := mocks.NewDifficultyLevelsInputPort(t)
		h := newTestDifficultyLevelsHandler(t, svc)

		id := uuid.MustParse("dddddddd-dddd-dddd-dddd-dddddddddddd")
		svc.EXPECT().Delete(mock.Anything, id).Return(nil).Once()

		r := httptest.NewRequest(http.MethodDelete, "/difficulty-levels/"+id.String(), nil)
		r = withRouteID(r, id.String())
		rr := httptest.NewRecorder()

		h.Delete(rr, r)

		require.Equal(t, http.StatusNoContent, rr.Code)
	})

	t.Run("service error maps to 500", func(t *testing.T) {
		svc := mocks.NewDifficultyLevelsInputPort(t)
		h := newTestDifficultyLevelsHandler(t, svc)

		id := uuid.MustParse("eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee")
		svc.EXPECT().Delete(mock.Anything, id).Return(errors.New("boom")).Once()

		r := httptest.NewRequest(http.MethodDelete, "/difficulty-levels/"+id.String(), nil)
		r = withRouteID(r, id.String())
		rr := httptest.NewRecorder()

		h.Delete(rr, r)

		require.Equal(t, http.StatusInternalServerError, rr.Code)
		errResp := decodeDifficultyLevelsErrorResponse(t, rr)
		require.Equal(t, utils.ErrorCodeInternalError, errResp.Error.Code)
	})
}
