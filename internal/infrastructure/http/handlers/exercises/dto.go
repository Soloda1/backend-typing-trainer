package exercises

import (
	"github.com/google/uuid"

	"backend-typing-trainer/internal/domain/models"
)

type upsertExerciseRequest struct {
	Text    string    `json:"text" validate:"required"`
	LevelID uuid.UUID `json:"level_id" validate:"required"`
}

type exerciseResponse struct {
	ID      uuid.UUID `json:"id"`
	Text    string    `json:"text"`
	LevelID uuid.UUID `json:"level_id"`
}

type exerciseSingleResponse struct {
	Exercise exerciseResponse `json:"exercise"`
}

type exerciseListResponse struct {
	Exercises []exerciseResponse `json:"exercises"`
}

func toExerciseResponse(exercise *models.Exercise) exerciseResponse {
	return exerciseResponse{
		ID:      exercise.ID,
		Text:    exercise.Text,
		LevelID: exercise.LevelID,
	}
}
