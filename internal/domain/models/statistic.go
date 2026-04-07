package models

import (
	"time"

	"github.com/google/uuid"
)

type Statistic struct {
	ID              uuid.UUID
	UserID          uuid.UUID
	LevelID         uuid.UUID
	ExerciseID      uuid.UUID
	MistakesPercent float64
	ExecutionTime   float64
	Speed           float64
	CreatedAt       time.Time
}
