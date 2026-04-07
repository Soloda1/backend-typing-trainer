package exercises

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/google/uuid"

	"backend-typing-trainer/internal/domain/models"
	exercisesport "backend-typing-trainer/internal/domain/ports/output/exercises"
	ports "backend-typing-trainer/internal/domain/ports/output/logger"
	"backend-typing-trainer/internal/utils"
)

type Service struct {
	repo exercisesport.Repository
	log  ports.Logger
}

func NewService(repo exercisesport.Repository, log ports.Logger) *Service {
	return &Service{
		repo: repo,
		log:  log.With(slog.String("component", "application_exercises")),
	}
}

func (s *Service) Create(ctx context.Context, exercise *models.Exercise) error {
	if err := validateExerciseForWrite(exercise, false); err != nil {
		s.log.Warn("create exercise rejected", slog.String("error", err.Error()))
		return err
	}

	exercise.Text = strings.TrimSpace(exercise.Text)
	if err := s.repo.Create(ctx, exercise); err != nil {
		if isKnownError(err) {
			s.log.Warn("create exercise failed", slog.String("error", err.Error()))
			return err
		}
		s.log.Error("create exercise failed", slog.String("error", err.Error()))
		return fmt.Errorf("create exercise: %w", err)
	}

	return nil
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*models.Exercise, error) {
	if id == uuid.Nil {
		s.log.Warn("get exercise by id rejected: empty id")
		return nil, utils.ErrInvalidRequest
	}

	exercise, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if isKnownError(err) {
			s.log.Warn("get exercise by id failed", slog.String("exercise_id", id.String()), slog.String("error", err.Error()))
			return nil, err
		}
		s.log.Error("get exercise by id failed", slog.String("exercise_id", id.String()), slog.String("error", err.Error()))
		return nil, fmt.Errorf("get exercise by id: %w", err)
	}

	return exercise, nil
}

func (s *Service) Update(ctx context.Context, exercise *models.Exercise) error {
	if err := validateExerciseForWrite(exercise, true); err != nil {
		s.log.Warn("update exercise rejected", slog.String("error", err.Error()))
		return err
	}

	exercise.Text = strings.TrimSpace(exercise.Text)
	if err := s.repo.Update(ctx, exercise); err != nil {
		if isKnownError(err) {
			s.log.Warn("update exercise failed", slog.String("exercise_id", exercise.ID.String()), slog.String("error", err.Error()))
			return err
		}
		s.log.Error("update exercise failed", slog.String("exercise_id", exercise.ID.String()), slog.String("error", err.Error()))
		return fmt.Errorf("update exercise: %w", err)
	}

	return nil
}

func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		s.log.Warn("delete exercise rejected: empty id")
		return utils.ErrInvalidRequest
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		if isKnownError(err) {
			s.log.Warn("delete exercise failed", slog.String("exercise_id", id.String()), slog.String("error", err.Error()))
			return err
		}
		s.log.Error("delete exercise failed", slog.String("exercise_id", id.String()), slog.String("error", err.Error()))
		return fmt.Errorf("delete exercise: %w", err)
	}

	return nil
}

func (s *Service) List(ctx context.Context, limit, offset int) ([]*models.Exercise, error) {
	if limit <= 0 || offset < 0 {
		s.log.Warn("list exercises rejected: invalid pagination", slog.Int("limit", limit), slog.Int("offset", offset))
		return nil, utils.ErrInvalidRequest
	}

	exercises, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		if isKnownError(err) {
			s.log.Warn("list exercises failed", slog.String("error", err.Error()))
			return nil, err
		}
		s.log.Error("list exercises failed", slog.String("error", err.Error()))
		return nil, fmt.Errorf("list exercises: %w", err)
	}

	return exercises, nil
}

func validateExerciseForWrite(exercise *models.Exercise, checkID bool) error {
	if exercise == nil {
		return utils.ErrInvalidRequest
	}
	if checkID && exercise.ID == uuid.Nil {
		return utils.ErrInvalidRequest
	}
	if strings.TrimSpace(exercise.Text) == "" {
		return utils.ErrInvalidRequest
	}
	if exercise.LevelID == uuid.Nil {
		return utils.ErrInvalidRequest
	}
	return nil
}

func isKnownError(err error) bool {
	return errors.Is(err, utils.ErrInvalidRequest) || errors.Is(err, utils.ErrNotFound)
}
