package keyboardzones

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/google/uuid"

	"backend-typing-trainer/internal/domain/models"
	keyboardzonesport "backend-typing-trainer/internal/domain/ports/output/keyboard_zones"
	ports "backend-typing-trainer/internal/domain/ports/output/logger"
	"backend-typing-trainer/internal/utils"
)

type Service struct {
	repo keyboardzonesport.Repository
	log  ports.Logger
}

func NewService(repo keyboardzonesport.Repository, log ports.Logger) *Service {
	return &Service{
		repo: repo,
		log:  log.With(slog.String("component", "application_keyboard_zones")),
	}
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*models.KeyboardZone, error) {
	if id == uuid.Nil {
		s.log.Warn("get keyboard zone by id rejected: empty id")
		return nil, utils.ErrInvalidRequest
	}

	zone, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if isKnownError(err) {
			s.log.Warn("get keyboard zone by id failed", slog.String("zone_id", id.String()), slog.String("error", err.Error()))
			return nil, err
		}
		s.log.Error("get keyboard zone by id failed", slog.String("zone_id", id.String()), slog.String("error", err.Error()))
		return nil, fmt.Errorf("get keyboard zone by id: %w", err)
	}

	return zone, nil
}

func (s *Service) GetByName(ctx context.Context, name string) (*models.KeyboardZone, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		s.log.Warn("get keyboard zone by name rejected: empty name")
		return nil, utils.ErrInvalidRequest
	}

	zone, err := s.repo.GetByName(ctx, name)
	if err != nil {
		if isKnownError(err) {
			s.log.Warn("get keyboard zone by name failed", slog.String("name", name), slog.String("error", err.Error()))
			return nil, err
		}
		s.log.Error("get keyboard zone by name failed", slog.String("name", name), slog.String("error", err.Error()))
		return nil, fmt.Errorf("get keyboard zone by name: %w", err)
	}

	return zone, nil
}

func (s *Service) List(ctx context.Context, limit, offset int) ([]*models.KeyboardZone, error) {
	if limit <= 0 || offset < 0 {
		s.log.Warn("list keyboard zones rejected: invalid pagination", slog.Int("limit", limit), slog.Int("offset", offset))
		return nil, utils.ErrInvalidRequest
	}

	zones, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		if isKnownError(err) {
			s.log.Warn("list keyboard zones failed", slog.String("error", err.Error()))
			return nil, err
		}
		s.log.Error("list keyboard zones failed", slog.String("error", err.Error()))
		return nil, fmt.Errorf("list keyboard zones: %w", err)
	}

	return zones, nil
}

func isKnownError(err error) bool {
	return errors.Is(err, utils.ErrInvalidRequest) || errors.Is(err, utils.ErrNotFound)
}
