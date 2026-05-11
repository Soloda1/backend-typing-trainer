package input

import (
	"context"

	"github.com/google/uuid"

	"backend-typing-trainer/internal/domain/models"
)

//go:generate mockery --name Users --structname UsersInputPort --dir . --output ../../../../mocks --outpkg mocks --with-expecter --filename UsersInputPort.go
type Users interface {
	GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetByLogin(ctx context.Context, login string) (*models.User, error)
	List(ctx context.Context, limit, offset int) ([]*models.User, error)
}
