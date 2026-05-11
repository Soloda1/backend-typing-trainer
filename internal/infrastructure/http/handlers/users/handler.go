package users

import (
	"backend-typing-trainer/internal/domain/ports/input"
	"backend-typing-trainer/internal/domain/ports/output/logger"
)

type Handler struct {
	usersService input.Users
	log          logger.Logger
}

func NewHandler(usersService input.Users, log logger.Logger) *Handler {
	return &Handler{
		usersService: usersService,
		log:          log,
	}
}
