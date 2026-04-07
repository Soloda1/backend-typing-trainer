package exercises_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	exercisesapp "backend-typing-trainer/internal/application/exercises"
	"backend-typing-trainer/internal/domain/models"
	"backend-typing-trainer/internal/utils"
	"backend-typing-trainer/mocks"
)

func newServiceWithMocks(t *testing.T, repo *mocks.ExercisesRepository) *exercisesapp.Service {
	t.Helper()

	log := mocks.NewLogger(t)
	log.EXPECT().With(mock.Anything).Return(log).Once()
	log.EXPECT().Warn(mock.Anything, mock.Anything).Maybe()
	log.EXPECT().Warn(mock.Anything, mock.Anything, mock.Anything).Maybe()
	log.EXPECT().Error(mock.Anything, mock.Anything).Maybe()
	log.EXPECT().Error(mock.Anything, mock.Anything, mock.Anything).Maybe()

	return exercisesapp.NewService(repo, log)
}

func validExercise() *models.Exercise {
	return &models.Exercise{
		Text:    "hello world",
		LevelID: uuid.MustParse("11111111-1111-1111-1111-111111111111"),
	}
}

func TestService_Create(t *testing.T) {
	repo := mocks.NewExercisesRepository(t)
	svc := newServiceWithMocks(t, repo)

	if err := svc.Create(context.Background(), nil); !errors.Is(err, utils.ErrInvalidRequest) {
		t.Fatalf("expected invalid request, got %v", err)
	}
	repo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)

	ex := validExercise()
	ex.Text = "   "
	if err := svc.Create(context.Background(), ex); !errors.Is(err, utils.ErrInvalidRequest) {
		t.Fatalf("expected invalid request, got %v", err)
	}

	repo.EXPECT().Create(mock.Anything, mock.AnythingOfType("*models.Exercise")).
		Run(func(_ context.Context, exercise *models.Exercise) {
			if exercise.Text != "hello world" {
				t.Fatalf("expected trimmed text, got %q", exercise.Text)
			}
		}).
		Return(nil).Once()
	if err := svc.Create(context.Background(), &models.Exercise{Text: "  hello world  ", LevelID: validExercise().LevelID}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestService_GetByID(t *testing.T) {
	repo := mocks.NewExercisesRepository(t)
	svc := newServiceWithMocks(t, repo)

	_, err := svc.GetByID(context.Background(), uuid.Nil)
	if !errors.Is(err, utils.ErrInvalidRequest) {
		t.Fatalf("expected invalid request, got %v", err)
	}

	id := uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa")
	repo.EXPECT().GetByID(mock.Anything, id).Return(&models.Exercise{ID: id}, nil).Once()
	_, err = svc.GetByID(context.Background(), id)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	repo.EXPECT().GetByID(mock.Anything, id).Return(nil, errors.New("db timeout")).Once()
	_, err = svc.GetByID(context.Background(), id)
	if err == nil {
		t.Fatal("expected wrapped repository error")
	}
}

func TestService_Update(t *testing.T) {
	repo := mocks.NewExercisesRepository(t)
	svc := newServiceWithMocks(t, repo)

	ex := validExercise()
	if err := svc.Update(context.Background(), ex); !errors.Is(err, utils.ErrInvalidRequest) {
		t.Fatalf("expected invalid request, got %v", err)
	}

	ex.ID = uuid.MustParse("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb")
	repo.EXPECT().Update(mock.Anything, ex).Return(utils.ErrNotFound).Once()
	if err := svc.Update(context.Background(), ex); !errors.Is(err, utils.ErrNotFound) {
		t.Fatalf("expected not found, got %v", err)
	}

	repo.EXPECT().Update(mock.Anything, ex).Return(errors.New("db timeout")).Once()
	if err := svc.Update(context.Background(), ex); err == nil {
		t.Fatal("expected wrapped repository error")
	}
}

func TestService_Delete(t *testing.T) {
	repo := mocks.NewExercisesRepository(t)
	svc := newServiceWithMocks(t, repo)

	id := uuid.MustParse("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb")

	if err := svc.Delete(context.Background(), uuid.Nil); !errors.Is(err, utils.ErrInvalidRequest) {
		t.Fatalf("expected invalid request, got %v", err)
	}

	repo.EXPECT().Delete(mock.Anything, id).Return(nil).Once()
	if err := svc.Delete(context.Background(), id); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	repo.EXPECT().Delete(mock.Anything, id).Return(errors.New("db timeout")).Once()
	if err := svc.Delete(context.Background(), id); err == nil {
		t.Fatal("expected wrapped repository error")
	}
}

func TestService_List(t *testing.T) {
	repo := mocks.NewExercisesRepository(t)
	svc := newServiceWithMocks(t, repo)

	id := uuid.MustParse("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb")

	_, err := svc.List(context.Background(), 0, 0)
	if !errors.Is(err, utils.ErrInvalidRequest) {
		t.Fatalf("expected invalid request, got %v", err)
	}

	repo.EXPECT().List(mock.Anything, 10, 0).Return([]*models.Exercise{{ID: id}}, nil).Once()
	list, err := svc.List(context.Background(), 10, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("expected 1 exercise, got %d", len(list))
	}

	repo.EXPECT().List(mock.Anything, 10, 0).Return(nil, errors.New("db timeout")).Once()
	_, err = svc.List(context.Background(), 10, 0)
	if err == nil {
		t.Fatal("expected wrapped repository error")
	}
}
