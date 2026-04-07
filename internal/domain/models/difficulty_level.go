package models

import "github.com/google/uuid"

type DifficultyLevel struct {
	ID                uuid.UUID
	AllowedMistakes   int
	KeyPressTime      float64
	MinExerciseLength int
	MaxExerciseLength int
	KeyboardZoneIDs   []uuid.UUID
}
