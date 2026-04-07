package models

import "github.com/google/uuid"

type KeyboardZone struct {
	ID      uuid.UUID
	Name    string
	Symbols string
}
