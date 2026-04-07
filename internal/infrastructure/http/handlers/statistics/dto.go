package statistics

import (
	"time"

	"github.com/google/uuid"

	"backend-typing-trainer/internal/domain/models"
)

type createStatisticRequest struct {
	UserID          uuid.UUID `json:"user_id" validate:"required"`
	LevelID         uuid.UUID `json:"level_id" validate:"required"`
	ExerciseID      uuid.UUID `json:"exercise_id" validate:"required"`
	MistakesPercent float64   `json:"mistakes_percent" validate:"gte=0"`
	ExecutionTime   float64   `json:"execution_time" validate:"gt=0"`
	Speed           float64   `json:"speed" validate:"gt=0"`
}

type statisticResponse struct {
	ID              uuid.UUID `json:"id"`
	UserID          uuid.UUID `json:"user_id"`
	LevelID         uuid.UUID `json:"level_id"`
	ExerciseID      uuid.UUID `json:"exercise_id"`
	MistakesPercent float64   `json:"mistakes_percent"`
	ExecutionTime   float64   `json:"execution_time"`
	Speed           float64   `json:"speed"`
	CreatedAt       time.Time `json:"created_at"`
}

type statisticSingleResponse struct {
	Statistic statisticResponse `json:"statistic"`
}

type statisticListResponse struct {
	Statistics []statisticResponse `json:"statistics"`
}

func toStatisticResponse(statistic *models.Statistic) statisticResponse {
	return statisticResponse{
		ID:              statistic.ID,
		UserID:          statistic.UserID,
		LevelID:         statistic.LevelID,
		ExerciseID:      statistic.ExerciseID,
		MistakesPercent: statistic.MistakesPercent,
		ExecutionTime:   statistic.ExecutionTime,
		Speed:           statistic.Speed,
		CreatedAt:       statistic.CreatedAt,
	}
}
