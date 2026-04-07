package models

import "github.com/google/uuid"

type LevelKeyboardZone struct {
	LevelID        uuid.UUID
	KeyboardZoneID uuid.UUID
}
