package auth

import (
	input "backend-typing-trainer/internal/domain/ports/input"
	ports "backend-typing-trainer/internal/domain/ports/output/logger"
)

type Handler struct {
	log         ports.Logger
	authService input.Auth
}

func NewHandler(log ports.Logger, authService input.Auth) *Handler {
	log = log.With("handler", "auth")

	return &Handler{
		log:         log,
		authService: authService,
	}
}
