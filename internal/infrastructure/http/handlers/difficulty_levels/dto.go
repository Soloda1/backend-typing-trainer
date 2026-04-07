package difficultylevels

import (
	"github.com/google/uuid"

	"backend-typing-trainer/internal/domain/models"
)

type upsertDifficultyLevelRequest struct {
	AllowedMistakes   int         `json:"allowed_mistakes" validate:"gte=0"`
	KeyPressTime      float64     `json:"key_press_time" validate:"gt=0"`
	MinExerciseLength int         `json:"min_exercise_length" validate:"gt=0"`
	MaxExerciseLength int         `json:"max_exercise_length" validate:"gt=0"`
	KeyboardZoneIDs   []uuid.UUID `json:"keyboard_zone_ids" validate:"required,min=1,dive,required"`
}

type difficultyLevelResponse struct {
	ID                uuid.UUID   `json:"id"`
	AllowedMistakes   int         `json:"allowed_mistakes"`
	KeyPressTime      float64     `json:"key_press_time"`
	MinExerciseLength int         `json:"min_exercise_length"`
	MaxExerciseLength int         `json:"max_exercise_length"`
	KeyboardZoneIDs   []uuid.UUID `json:"keyboard_zone_ids"`
}

type difficultyLevelSingleResponse struct {
	DifficultyLevel difficultyLevelResponse `json:"difficulty_level"`
}

type difficultyLevelListResponse struct {
	DifficultyLevels []difficultyLevelResponse `json:"difficulty_levels"`
}

func toDifficultyLevelResponse(level *models.DifficultyLevel) difficultyLevelResponse {
	return difficultyLevelResponse{
		ID:                level.ID,
		AllowedMistakes:   level.AllowedMistakes,
		KeyPressTime:      level.KeyPressTime,
		MinExerciseLength: level.MinExerciseLength,
		MaxExerciseLength: level.MaxExerciseLength,
		KeyboardZoneIDs:   level.KeyboardZoneIDs,
	}
}
