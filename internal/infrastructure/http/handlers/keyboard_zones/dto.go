package keyboardzones

import (
	"github.com/google/uuid"

	"backend-typing-trainer/internal/domain/models"
)

type keyboardZoneResponse struct {
	ID      uuid.UUID `json:"id"`
	Name    string    `json:"name"`
	Symbols string    `json:"symbols"`
}

type keyboardZoneSingleResponse struct {
	KeyboardZone keyboardZoneResponse `json:"keyboard_zone"`
}

type keyboardZoneListResponse struct {
	KeyboardZones []keyboardZoneResponse `json:"keyboard_zones"`
}

func toKeyboardZoneResponse(zone *models.KeyboardZone) keyboardZoneResponse {
	return keyboardZoneResponse{
		ID:      zone.ID,
		Name:    zone.Name,
		Symbols: zone.Symbols,
	}
}
