package difficultylevels

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/google/uuid"

	"backend-typing-trainer/internal/domain/models"
	difficultylevelsport "backend-typing-trainer/internal/domain/ports/output/difficulty_levels"
	ports "backend-typing-trainer/internal/domain/ports/output/logger"
	"backend-typing-trainer/internal/utils"
)

type Service struct {
	repo difficultylevelsport.Repository
	log  ports.Logger
}

func NewService(repo difficultylevelsport.Repository, log ports.Logger) *Service {
	return &Service{
		repo: repo,
		log:  log.With(slog.String("component", "application_difficulty_levels")),
	}
}

func (s *Service) Create(ctx context.Context, level *models.DifficultyLevel) error {
	if err := validateDifficultyLevelForWrite(level, false); err != nil {
		s.log.Warn("create difficulty level rejected", slog.String("error", err.Error()))
		return err
	}

	if err := s.repo.Create(ctx, level); err != nil {
		if isKnownError(err) {
			s.log.Warn("create difficulty level failed", slog.String("error", err.Error()))
			return err
		}
		s.log.Error("create difficulty level failed", slog.String("error", err.Error()))
		return fmt.Errorf("create difficulty level: %w", err)
	}

	return nil
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*models.DifficultyLevel, error) {
	if id == uuid.Nil {
		s.log.Warn("get difficulty level rejected: empty id")
		return nil, utils.ErrInvalidRequest
	}

	level, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if isKnownError(err) {
			s.log.Warn("get difficulty level failed", slog.String("level_id", id.String()), slog.String("error", err.Error()))
			return nil, err
		}
		s.log.Error("get difficulty level failed", slog.String("level_id", id.String()), slog.String("error", err.Error()))
		return nil, fmt.Errorf("get difficulty level by id: %w", err)
	}

	return level, nil
}

func (s *Service) Update(ctx context.Context, level *models.DifficultyLevel) error {
	if err := validateDifficultyLevelForWrite(level, true); err != nil {
		s.log.Warn("update difficulty level rejected", slog.String("error", err.Error()))
		return err
	}

	if err := s.repo.Update(ctx, level); err != nil {
		if isKnownError(err) {
			s.log.Warn("update difficulty level failed", slog.String("level_id", level.ID.String()), slog.String("error", err.Error()))
			return err
		}
		s.log.Error("update difficulty level failed", slog.String("level_id", level.ID.String()), slog.String("error", err.Error()))
		return fmt.Errorf("update difficulty level: %w", err)
	}

	return nil
}

func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		s.log.Warn("delete difficulty level rejected: empty id")
		return utils.ErrInvalidRequest
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		if isKnownError(err) {
			s.log.Warn("delete difficulty level failed", slog.String("level_id", id.String()), slog.String("error", err.Error()))
			return err
		}
		s.log.Error("delete difficulty level failed", slog.String("level_id", id.String()), slog.String("error", err.Error()))
		return fmt.Errorf("delete difficulty level: %w", err)
	}

	return nil
}

func (s *Service) List(ctx context.Context, limit, offset int) ([]*models.DifficultyLevel, error) {
	if limit <= 0 || offset < 0 {
		s.log.Warn("list difficulty levels rejected: invalid pagination", slog.Int("limit", limit), slog.Int("offset", offset))
		return nil, utils.ErrInvalidRequest
	}

	levels, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		if isKnownError(err) {
			s.log.Warn("list difficulty levels failed", slog.String("error", err.Error()))
			return nil, err
		}
		s.log.Error("list difficulty levels failed", slog.String("error", err.Error()))
		return nil, fmt.Errorf("list difficulty levels: %w", err)
	}

	return levels, nil
}

func validateDifficultyLevelForWrite(level *models.DifficultyLevel, checkID bool) error {
	if level == nil {
		return utils.ErrInvalidRequest
	}
	if checkID && level.ID == uuid.Nil {
		return utils.ErrInvalidRequest
	}
	if level.AllowedMistakes < 0 || level.KeyPressTime <= 0 {
		return utils.ErrInvalidRequest
	}
	if level.MinExerciseLength <= 0 || level.MaxExerciseLength <= 0 || level.MinExerciseLength > level.MaxExerciseLength {
		return utils.ErrInvalidRequest
	}
	if len(level.KeyboardZoneIDs) == 0 {
		return utils.ErrInvalidRequest
	}

	seen := make(map[uuid.UUID]struct{}, len(level.KeyboardZoneIDs))
	for _, zoneID := range level.KeyboardZoneIDs {
		if zoneID == uuid.Nil {
			return utils.ErrInvalidRequest
		}
		if _, exists := seen[zoneID]; exists {
			return utils.ErrInvalidRequest
		}
		seen[zoneID] = struct{}{}
	}

	return nil
}

func isKnownError(err error) bool {
	return errors.Is(err, utils.ErrInvalidRequest) || errors.Is(err, utils.ErrNotFound)
}
