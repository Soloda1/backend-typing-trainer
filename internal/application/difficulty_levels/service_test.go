package difficultylevels_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	difficultylevelsapp "backend-typing-trainer/internal/application/difficulty_levels"
	"backend-typing-trainer/internal/domain/models"
	"backend-typing-trainer/internal/utils"
	"backend-typing-trainer/mocks"
)

func newServiceWithLoggerMock(t *testing.T, repo *mocks.DifficultyLevelsRepository) *difficultylevelsapp.Service {
	t.Helper()

	log := mocks.NewLogger(t)
	log.EXPECT().With(mock.Anything).Return(log).Once()
	log.EXPECT().Warn(mock.Anything, mock.Anything).Maybe()
	log.EXPECT().Warn(mock.Anything, mock.Anything, mock.Anything).Maybe()
	log.EXPECT().Error(mock.Anything, mock.Anything).Maybe()
	log.EXPECT().Error(mock.Anything, mock.Anything, mock.Anything).Maybe()

	return difficultylevelsapp.NewService(repo, log)
}

func validLevel() *models.DifficultyLevel {
	return &models.DifficultyLevel{
		AllowedMistakes:   2,
		KeyPressTime:      1.5,
		MinExerciseLength: 5,
		MaxExerciseLength: 25,
		KeyboardZoneIDs: []uuid.UUID{
			uuid.MustParse("11111111-1111-1111-1111-111111111111"),
			uuid.MustParse("22222222-2222-2222-2222-222222222222"),
		},
	}
}

func TestServiceCreate(t *testing.T) {
	repo := mocks.NewDifficultyLevelsRepository(t)
	svc := newServiceWithLoggerMock(t, repo)

	if err := svc.Create(context.Background(), nil); !errors.Is(err, utils.ErrInvalidRequest) {
		t.Fatalf("expected invalid request for nil level, got %v", err)
	}
	repo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)

	invalid := validLevel()
	invalid.KeyboardZoneIDs = nil
	if err := svc.Create(context.Background(), invalid); !errors.Is(err, utils.ErrInvalidRequest) {
		t.Fatalf("expected invalid request for empty zones, got %v", err)
	}
	repo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)

	repo.EXPECT().Create(mock.Anything, mock.Anything).Return(errors.New("db timeout")).Once()
	if err := svc.Create(context.Background(), validLevel()); err == nil {
		t.Fatal("expected wrapped repository error")
	}
	repo.AssertNumberOfCalls(t, "Create", 1)
}

func TestServiceGetByID(t *testing.T) {
	repo := mocks.NewDifficultyLevelsRepository(t)
	svc := newServiceWithLoggerMock(t, repo)

	_, err := svc.GetByID(context.Background(), uuid.Nil)
	if !errors.Is(err, utils.ErrInvalidRequest) {
		t.Fatalf("expected invalid request, got %v", err)
	}
	repo.AssertNotCalled(t, "GetByID", mock.Anything, mock.Anything)

	id := uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa")
	expected := validLevel()
	expected.ID = id
	repo.EXPECT().GetByID(mock.Anything, id).Return(expected, nil).Once()

	got, err := svc.GetByID(context.Background(), id)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got == nil || got.ID != id {
		t.Fatalf("expected level with id %s", id)
	}
}

func TestServiceUpdateAndDelete(t *testing.T) {
	repo := mocks.NewDifficultyLevelsRepository(t)
	svc := newServiceWithLoggerMock(t, repo)

	lvl := validLevel()
	if err := svc.Update(context.Background(), lvl); !errors.Is(err, utils.ErrInvalidRequest) {
		t.Fatalf("expected invalid request for missing id, got %v", err)
	}
	repo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)

	lvl.ID = uuid.MustParse("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb")
	repo.EXPECT().Update(mock.Anything, lvl).Return(utils.ErrNotFound).Once()
	if err := svc.Update(context.Background(), lvl); !errors.Is(err, utils.ErrNotFound) {
		t.Fatalf("expected not found, got %v", err)
	}
	repo.AssertNumberOfCalls(t, "Update", 1)

	if err := svc.Delete(context.Background(), uuid.Nil); !errors.Is(err, utils.ErrInvalidRequest) {
		t.Fatalf("expected invalid request for empty id, got %v", err)
	}
	repo.AssertNotCalled(t, "Delete", mock.Anything, mock.Anything)

	repo.EXPECT().Delete(mock.Anything, lvl.ID).Return(nil).Once()
	if err := svc.Delete(context.Background(), lvl.ID); err != nil {
		t.Fatalf("unexpected delete error: %v", err)
	}
	repo.AssertNumberOfCalls(t, "Delete", 1)
}

func TestServiceList(t *testing.T) {
	repo := mocks.NewDifficultyLevelsRepository(t)
	svc := newServiceWithLoggerMock(t, repo)

	_, err := svc.List(context.Background(), 0, 0)
	if !errors.Is(err, utils.ErrInvalidRequest) {
		t.Fatalf("expected invalid request for limit, got %v", err)
	}
	repo.AssertNotCalled(t, "List", mock.Anything, mock.Anything, mock.Anything)

	expected := []*models.DifficultyLevel{{ID: uuid.MustParse("cccccccc-cccc-cccc-cccc-cccccccccccc")}}
	repo.EXPECT().List(mock.Anything, 10, 0).Return(expected, nil).Once()

	got, err := svc.List(context.Background(), 10, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 || got[0].ID != expected[0].ID {
		t.Fatalf("unexpected list result: %+v", got)
	}
	repo.AssertNumberOfCalls(t, "List", 1)
}
