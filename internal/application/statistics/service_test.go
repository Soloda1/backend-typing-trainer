package statistics_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	statisticsapp "backend-typing-trainer/internal/application/statistics"
	"backend-typing-trainer/internal/domain/models"
	"backend-typing-trainer/internal/utils"
	"backend-typing-trainer/mocks"
)

func newServiceWithMocks(t *testing.T, repo *mocks.StatisticsRepository) *statisticsapp.Service {
	t.Helper()

	log := mocks.NewLogger(t)
	log.EXPECT().With(mock.Anything).Return(log).Once()
	log.EXPECT().Warn(mock.Anything, mock.Anything).Maybe()
	log.EXPECT().Warn(mock.Anything, mock.Anything, mock.Anything).Maybe()
	log.EXPECT().Error(mock.Anything, mock.Anything).Maybe()
	log.EXPECT().Error(mock.Anything, mock.Anything, mock.Anything).Maybe()

	return statisticsapp.NewService(repo, log)
}

func validStatistic() *models.Statistic {
	return &models.Statistic{
		UserID:          uuid.MustParse("11111111-1111-1111-1111-111111111111"),
		LevelID:         uuid.MustParse("22222222-2222-2222-2222-222222222222"),
		ExerciseID:      uuid.MustParse("33333333-3333-3333-3333-333333333333"),
		MistakesPercent: 5,
		ExecutionTime:   10,
		Speed:           250,
	}
}

func TestService_Create(t *testing.T) {
	repo := mocks.NewStatisticsRepository(t)
	svc := newServiceWithMocks(t, repo)

	if err := svc.Create(context.Background(), nil); !errors.Is(err, utils.ErrInvalidRequest) {
		t.Fatalf("expected invalid request, got %v", err)
	}

	stat := validStatistic()
	repo.EXPECT().Create(mock.Anything, stat).Return(nil).Once()
	if err := svc.Create(context.Background(), stat); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestService_ListByUserID(t *testing.T) {
	repo := mocks.NewStatisticsRepository(t)
	svc := newServiceWithMocks(t, repo)

	_, err := svc.ListByUserID(context.Background(), uuid.Nil, 10, 0)
	if !errors.Is(err, utils.ErrInvalidRequest) {
		t.Fatalf("expected invalid request, got %v", err)
	}

	userID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	repo.EXPECT().ListByUserID(mock.Anything, userID, 10, 0).Return([]*models.Statistic{}, nil).Once()
	_, err = svc.ListByUserID(context.Background(), userID, 10, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestService_ListByLevelID(t *testing.T) {
	repo := mocks.NewStatisticsRepository(t)
	svc := newServiceWithMocks(t, repo)

	levelID := uuid.MustParse("22222222-2222-2222-2222-222222222222")
	repo.EXPECT().ListByLevelID(mock.Anything, levelID, 10, 0).Return(nil, errors.New("db timeout")).Once()
	_, err := svc.ListByLevelID(context.Background(), levelID, 10, 0)
	if err == nil {
		t.Fatal("expected wrapped repository error")
	}
}

func TestService_ListByExerciseID(t *testing.T) {
	repo := mocks.NewStatisticsRepository(t)
	svc := newServiceWithMocks(t, repo)

	exerciseID := uuid.MustParse("33333333-3333-3333-3333-333333333333")
	repo.EXPECT().ListByExerciseID(mock.Anything, exerciseID, 10, 0).Return([]*models.Statistic{{ExerciseID: exerciseID}}, nil).Once()
	stats, err := svc.ListByExerciseID(context.Background(), exerciseID, 10, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(stats) != 1 {
		t.Fatalf("expected 1 item, got %d", len(stats))
	}
}
