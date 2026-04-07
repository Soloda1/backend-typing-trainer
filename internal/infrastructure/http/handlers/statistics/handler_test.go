package statistics

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"backend-typing-trainer/internal/domain/models"
	infraLogger "backend-typing-trainer/internal/infrastructure/logger"
	"backend-typing-trainer/internal/utils"
	"backend-typing-trainer/mocks"
)

func newTestStatisticsHandler(t *testing.T, svc *mocks.StatisticsInputPort) *Handler {
	t.Helper()
	return NewHandler(infraLogger.New("dev"), svc)
}

func decodeStatisticsError(t *testing.T, rr *httptest.ResponseRecorder) utils.ErrorResponse {
	t.Helper()
	var resp utils.ErrorResponse
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &resp))
	return resp
}

func TestHandler_Create(t *testing.T) {
	svc := mocks.NewStatisticsInputPort(t)
	h := newTestStatisticsHandler(t, svc)

	svc.EXPECT().Create(mock.Anything, mock.AnythingOfType("*models.Statistic")).Return(nil).Once()

	r := httptest.NewRequest(http.MethodPost, "/statistics", strings.NewReader(`{"user_id":"11111111-1111-1111-1111-111111111111","level_id":"22222222-2222-2222-2222-222222222222","exercise_id":"33333333-3333-3333-3333-333333333333","mistakes_percent":5,"execution_time":10,"speed":250}`))
	rr := httptest.NewRecorder()

	h.Create(rr, r)

	require.Equal(t, http.StatusCreated, rr.Code)
}

func TestHandler_ListByUser(t *testing.T) {
	svc := mocks.NewStatisticsInputPort(t)
	h := newTestStatisticsHandler(t, svc)

	userID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	svc.EXPECT().ListByUserID(mock.Anything, userID, 20, 0).Return([]*models.Statistic{}, nil).Once()

	r := httptest.NewRequest(http.MethodGet, "/statistics?user_id="+userID.String(), nil)
	rr := httptest.NewRecorder()

	h.List(rr, r)

	require.Equal(t, http.StatusOK, rr.Code)
}

func TestHandler_ListByLevel(t *testing.T) {
	svc := mocks.NewStatisticsInputPort(t)
	h := newTestStatisticsHandler(t, svc)

	levelID := uuid.MustParse("22222222-2222-2222-2222-222222222222")
	svc.EXPECT().ListByLevelID(mock.Anything, levelID, 20, 0).Return([]*models.Statistic{}, nil).Once()

	r := httptest.NewRequest(http.MethodGet, "/statistics?level_id="+levelID.String(), nil)
	rr := httptest.NewRecorder()

	h.List(rr, r)

	require.Equal(t, http.StatusOK, rr.Code)
}

func TestHandler_ListByExercise(t *testing.T) {
	svc := mocks.NewStatisticsInputPort(t)
	h := newTestStatisticsHandler(t, svc)

	exerciseID := uuid.MustParse("33333333-3333-3333-3333-333333333333")
	svc.EXPECT().ListByExerciseID(mock.Anything, exerciseID, 20, 0).Return([]*models.Statistic{}, nil).Once()

	r := httptest.NewRequest(http.MethodGet, "/statistics?exercise_id="+exerciseID.String(), nil)
	rr := httptest.NewRecorder()

	h.List(rr, r)

	require.Equal(t, http.StatusOK, rr.Code)
}

func TestHandler_ListInvalidFilters(t *testing.T) {
	svc := mocks.NewStatisticsInputPort(t)
	h := newTestStatisticsHandler(t, svc)

	r := httptest.NewRequest(http.MethodGet, "/statistics", nil)
	rr := httptest.NewRecorder()

	h.List(rr, r)

	require.Equal(t, http.StatusBadRequest, rr.Code)
	errResp := decodeStatisticsError(t, rr)
	require.Equal(t, utils.ErrorCodeInvalidRequest, errResp.Error.Code)
}
