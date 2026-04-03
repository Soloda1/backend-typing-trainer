package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID
	Login        string
	Role         UserRole
	PasswordHash *string
	CreatedAt    time.Time
}
