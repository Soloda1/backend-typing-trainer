package models

import "github.com/google/uuid"

type Exercise struct {
	ID      uuid.UUID
	Text    string
	LevelID uuid.UUID
}
