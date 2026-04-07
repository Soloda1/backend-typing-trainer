package difficultylevels

import (
	"context"

	"github.com/google/uuid"

	"backend-typing-trainer/internal/domain/models"
)

//go:generate mockery --name Repository --structname DifficultyLevelsRepository --dir . --output ../../../../../mocks --outpkg mocks --with-expecter --filename DifficultyLevelsRepository.go

type Repository interface {
	Create(ctx context.Context, level *models.DifficultyLevel) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.DifficultyLevel, error)
	Update(ctx context.Context, level *models.DifficultyLevel) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int) ([]*models.DifficultyLevel, error)
}
