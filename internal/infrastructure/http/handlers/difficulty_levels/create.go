package difficultylevels

import (
	"log/slog"
	"net/http"

	"backend-typing-trainer/internal/domain/models"
	"backend-typing-trainer/internal/utils"
)

// Create godoc
// @Summary Создание уровня сложности
// @Description Доступно только администратору. Создает уровень сложности и связывает его с зонами клавиатуры.
// @Tags DifficultyLevels
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body upsertDifficultyLevelRequest true "Difficulty level payload"
// @Success 201 {object} difficultyLevelSingleResponse "Уровень сложности создан"
// @Failure 400 {object} utils.ErrorResponse "Неверный запрос"
// @Failure 401 {object} utils.ErrorResponse "Не авторизован"
// @Failure 403 {object} utils.ErrorResponse "Недостаточно прав"
// @Failure 500 {object} utils.InternalErrorResponse "Внутренняя ошибка сервера"
// @Router /difficulty-levels [post]
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	h.log.Debug("create difficulty level started")

	if h.difficultyLevelsService == nil {
		h.log.Error("create difficulty level failed: service is nil")
		if err := utils.WriteError(w, http.StatusInternalServerError, utils.ErrorCodeInternalError, "internal server error"); err != nil {
			h.log.Error("create difficulty level: write 500 response failed", slog.String("error", err.Error()))
		}
		return
	}

	var req upsertDifficultyLevelRequest
	if ok := h.decodeAndValidate(w, r, &req); !ok {
		return
	}

	level := &models.DifficultyLevel{
		AllowedMistakes:   req.AllowedMistakes,
		KeyPressTime:      req.KeyPressTime,
		MinExerciseLength: req.MinExerciseLength,
		MaxExerciseLength: req.MaxExerciseLength,
		KeyboardZoneIDs:   req.KeyboardZoneIDs,
	}

	if err := h.difficultyLevelsService.Create(r.Context(), level); err != nil {
		h.log.Warn("create difficulty level failed", slog.String("error", err.Error()))
		apiErr := utils.MapError(err)
		if writeErr := utils.WriteError(w, apiErr.Status, apiErr.Code, apiErr.Message); writeErr != nil {
			h.log.Error("create difficulty level: write service error response failed", slog.Int("status", apiErr.Status), slog.String("code", string(apiErr.Code)), slog.String("error", writeErr.Error()))
		}
		return
	}

	if err := utils.WriteJSON(w, http.StatusCreated, difficultyLevelSingleResponse{DifficultyLevel: toDifficultyLevelResponse(level)}); err != nil {
		h.log.Error("create difficulty level: write success response failed", slog.String("error", err.Error()))
	}
}
