package input

import (
	"context"

	"github.com/google/uuid"

	"backend-typing-trainer/internal/domain/models"
)

//go:generate mockery --name Exercises --structname ExercisesInputPort --dir . --output ../../../../mocks --outpkg mocks --with-expecter --filename ExercisesInputPort.go

type Exercises interface {
	Create(ctx context.Context, exercise *models.Exercise) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Exercise, error)
	Update(ctx context.Context, exercise *models.Exercise) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int) ([]*models.Exercise, error)
}
