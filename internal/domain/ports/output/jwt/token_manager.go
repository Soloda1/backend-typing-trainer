package jwt

import (
	"backend-typing-trainer/internal/domain/models"

	"github.com/google/uuid"
)

type TokenClaims struct {
	UserID uuid.UUID
	Role   models.UserRole
}

//go:generate mockery --name TokenManager --dir . --output ../../../../../mocks --outpkg mocks --with-expecter --filename TokenManager.go

type TokenManager interface {
	NewToken(userID uuid.UUID, role models.UserRole) (string, error)
	ParseToken(token string) (*TokenClaims, error)
}
