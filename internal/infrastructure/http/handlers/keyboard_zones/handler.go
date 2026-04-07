package keyboardzones

import (
	input "backend-typing-trainer/internal/domain/ports/input"
	ports "backend-typing-trainer/internal/domain/ports/output/logger"
)

type Handler struct {
	log                  ports.Logger
	keyboardZonesService input.KeyboardZones
}

func NewHandler(log ports.Logger, keyboardZonesService input.KeyboardZones) *Handler {
	log = log.With("handler", "keyboard_zones")

	return &Handler{
		log:                  log,
		keyboardZonesService: keyboardZonesService,
	}
}
