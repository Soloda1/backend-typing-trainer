package users

import (
	"context"
	"log/slog"
	"strings"

	"github.com/google/uuid"

	"backend-typing-trainer/internal/domain/models"
	"backend-typing-trainer/internal/domain/ports/output/logger"
	"backend-typing-trainer/internal/domain/ports/output/users"
	"backend-typing-trainer/internal/utils"
)

type Service struct {
	repo users.Repository
	log  logger.Logger
}

func NewService(repo users.Repository, log logger.Logger) *Service {
	return &Service{
		repo: repo,
		log:  log.With(slog.String("component", "application_users")),
	}
}
func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	if id == uuid.Nil {
		return nil, utils.ErrInvalidRequest
	}
	return s.repo.GetByID(ctx, id)
}
func (s *Service) GetByLogin(ctx context.Context, login string) (*models.User, error) {
	login = strings.TrimSpace(login)
	if login == "" {
		return nil, utils.ErrInvalidRequest
	}
	return s.repo.GetByLogin(ctx, login)
}
func (s *Service) List(ctx context.Context, limit, offset int) ([]*models.User, error) {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}
	return s.repo.List(ctx, limit, offset)
}
