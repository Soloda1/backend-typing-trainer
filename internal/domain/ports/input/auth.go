package input

import (
	"context"

	"backend-typing-trainer/internal/domain/models"
)

//go:generate mockery --name Auth --structname AuthInputPort --dir . --output ../../../../mocks --outpkg mocks --with-expecter --filename AuthInputPort.go

type Auth interface {
	Register(ctx context.Context, login, password string, role models.UserRole) (*models.User, error)
	Login(ctx context.Context, login, password string) (token string, err error)
}
