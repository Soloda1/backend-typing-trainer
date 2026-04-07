package statistics

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/google/uuid"

	"backend-typing-trainer/internal/domain/models"
	ports "backend-typing-trainer/internal/domain/ports/output/logger"
	statisticsport "backend-typing-trainer/internal/domain/ports/output/statistics"
	"backend-typing-trainer/internal/utils"
)

type Service struct {
	repo statisticsport.Repository
	log  ports.Logger
}

func NewService(repo statisticsport.Repository, log ports.Logger) *Service {
	return &Service{
		repo: repo,
		log:  log.With(slog.String("component", "application_statistics")),
	}
}

func (s *Service) Create(ctx context.Context, statistic *models.Statistic) error {
	if err := validateStatisticForCreate(statistic); err != nil {
		s.log.Warn("create statistic rejected", slog.String("error", err.Error()))
		return err
	}

	if err := s.repo.Create(ctx, statistic); err != nil {
		if isKnownError(err) {
			s.log.Warn("create statistic failed", slog.String("error", err.Error()))
			return err
		}
		s.log.Error("create statistic failed", slog.String("error", err.Error()))
		return fmt.Errorf("create statistic: %w", err)
	}

	return nil
}

func (s *Service) ListByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*models.Statistic, error) {
	if err := validateFilterAndPagination(userID, limit, offset); err != nil {
		s.log.Warn("list statistics by user rejected", slog.String("error", err.Error()))
		return nil, err
	}

	statistics, err := s.repo.ListByUserID(ctx, userID, limit, offset)
	if err != nil {
		if isKnownError(err) {
			s.log.Warn("list statistics by user failed", slog.String("error", err.Error()))
			return nil, err
		}
		s.log.Error("list statistics by user failed", slog.String("error", err.Error()))
		return nil, fmt.Errorf("list statistics by user: %w", err)
	}

	return statistics, nil
}

func (s *Service) ListByLevelID(ctx context.Context, levelID uuid.UUID, limit, offset int) ([]*models.Statistic, error) {
	if err := validateFilterAndPagination(levelID, limit, offset); err != nil {
		s.log.Warn("list statistics by level rejected", slog.String("error", err.Error()))
		return nil, err
	}

	statistics, err := s.repo.ListByLevelID(ctx, levelID, limit, offset)
	if err != nil {
		if isKnownError(err) {
			s.log.Warn("list statistics by level failed", slog.String("error", err.Error()))
			return nil, err
		}
		s.log.Error("list statistics by level failed", slog.String("error", err.Error()))
		return nil, fmt.Errorf("list statistics by level: %w", err)
	}

	return statistics, nil
}

func (s *Service) ListByExerciseID(ctx context.Context, exerciseID uuid.UUID, limit, offset int) ([]*models.Statistic, error) {
	if err := validateFilterAndPagination(exerciseID, limit, offset); err != nil {
		s.log.Warn("list statistics by exercise rejected", slog.String("error", err.Error()))
		return nil, err
	}

	statistics, err := s.repo.ListByExerciseID(ctx, exerciseID, limit, offset)
	if err != nil {
		if isKnownError(err) {
			s.log.Warn("list statistics by exercise failed", slog.String("error", err.Error()))
			return nil, err
		}
		s.log.Error("list statistics by exercise failed", slog.String("error", err.Error()))
		return nil, fmt.Errorf("list statistics by exercise: %w", err)
	}

	return statistics, nil
}

func validateStatisticForCreate(statistic *models.Statistic) error {
	if statistic == nil {
		return utils.ErrInvalidRequest
	}
	if statistic.UserID == uuid.Nil || statistic.LevelID == uuid.Nil || statistic.ExerciseID == uuid.Nil {
		return utils.ErrInvalidRequest
	}
	if statistic.MistakesPercent < 0 || statistic.ExecutionTime <= 0 || statistic.Speed <= 0 {
		return utils.ErrInvalidRequest
	}
	return nil
}

func validateFilterAndPagination(id uuid.UUID, limit, offset int) error {
	if id == uuid.Nil || limit <= 0 || offset < 0 {
		return utils.ErrInvalidRequest
	}
	return nil
}

func isKnownError(err error) bool {
	return errors.Is(err, utils.ErrInvalidRequest) || errors.Is(err, utils.ErrNotFound)
}
