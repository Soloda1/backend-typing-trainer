package exercises

import (
	"log/slog"
	"net/http"

	"backend-typing-trainer/internal/domain/models"
	"backend-typing-trainer/internal/utils"
)

// Create godoc
// @Summary Создание упражнения
// @Description Доступно только администратору.
// @Tags Exercises
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body upsertExerciseRequest true "Exercise payload"
// @Success 201 {object} exerciseSingleResponse "Упражнение создано"
// @Failure 400 {object} utils.ErrorResponse "Неверный запрос"
// @Failure 401 {object} utils.ErrorResponse "Не авторизован"
// @Failure 403 {object} utils.ErrorResponse "Недостаточно прав"
// @Failure 500 {object} utils.InternalErrorResponse "Внутренняя ошибка сервера"
// @Router /exercises [post]
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	h.log.Debug("create exercise started")

	if h.exercisesService == nil {
		h.log.Error("create exercise failed: service is nil")
		if err := utils.WriteError(w, http.StatusInternalServerError, utils.ErrorCodeInternalError, "internal server error"); err != nil {
			h.log.Error("create exercise: write 500 response failed", slog.String("error", err.Error()))
		}
		return
	}

	var req upsertExerciseRequest
	if ok := h.decodeAndValidate(w, r, &req); !ok {
		return
	}

	exercise := &models.Exercise{
		Text:    req.Text,
		LevelID: req.LevelID,
	}

	if err := h.exercisesService.Create(r.Context(), exercise); err != nil {
		h.log.Warn("create exercise failed", slog.String("error", err.Error()))
		apiErr := utils.MapError(err)
		if writeErr := utils.WriteError(w, apiErr.Status, apiErr.Code, apiErr.Message); writeErr != nil {
			h.log.Error("create exercise: write service error response failed", slog.Int("status", apiErr.Status), slog.String("code", string(apiErr.Code)), slog.String("error", writeErr.Error()))
		}
		return
	}

	if err := utils.WriteJSON(w, http.StatusCreated, exerciseSingleResponse{Exercise: toExerciseResponse(exercise)}); err != nil {
		h.log.Error("create exercise: write success response failed", slog.String("error", err.Error()))
	}
}
