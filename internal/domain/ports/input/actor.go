package input

import (
	"github.com/google/uuid"

	"backend-typing-trainer/internal/domain/models"
)

type Actor struct {
	UserID uuid.UUID
	Role   models.UserRole
}
