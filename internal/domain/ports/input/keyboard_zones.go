package input

import (
	"context"

	"github.com/google/uuid"

	"backend-typing-trainer/internal/domain/models"
)

//go:generate mockery --name KeyboardZones --structname KeyboardZonesInputPort --dir . --output ../../../../mocks --outpkg mocks --with-expecter --filename KeyboardZonesInputPort.go

type KeyboardZones interface {
	GetByID(ctx context.Context, id uuid.UUID) (*models.KeyboardZone, error)
	GetByName(ctx context.Context, name string) (*models.KeyboardZone, error)
	List(ctx context.Context, limit, offset int) ([]*models.KeyboardZone, error)
}
