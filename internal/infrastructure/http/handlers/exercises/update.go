package exercises

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"backend-typing-trainer/internal/domain/models"
	"backend-typing-trainer/internal/utils"
)

// Update godoc
// @Summary Обновление упражнения
// @Description Доступно только администратору.
// @Tags Exercises
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Exercise ID" format(uuid)
// @Param request body upsertExerciseRequest true "Exercise payload"
// @Success 200 {object} exerciseSingleResponse "Упражнение обновлено"
// @Failure 400 {object} utils.ErrorResponse "Неверный запрос"
// @Failure 401 {object} utils.ErrorResponse "Не авторизован"
// @Failure 403 {object} utils.ErrorResponse "Недостаточно прав"
// @Failure 404 {object} utils.ErrorResponse "Упражнение не найдено"
// @Failure 500 {object} utils.InternalErrorResponse "Внутренняя ошибка сервера"
// @Router /exercises/{id} [patch]
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	h.log.Debug("update exercise started")

	if h.exercisesService == nil {
		h.log.Error("update exercise failed: service is nil")
		if err := utils.WriteError(w, http.StatusInternalServerError, utils.ErrorCodeInternalError, "internal server error"); err != nil {
			h.log.Error("update exercise: write 500 response failed", slog.String("error", err.Error()))
		}
		return
	}

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		h.log.Warn("update exercise failed: invalid id", slog.String("error", err.Error()))
		apiErr := utils.MapError(utils.ErrInvalidRequest)
		if writeErr := utils.WriteError(w, apiErr.Status, apiErr.Code, apiErr.Message); writeErr != nil {
			h.log.Error("update exercise: write invalid id error response failed", slog.Int("status", apiErr.Status), slog.String("code", string(apiErr.Code)), slog.String("error", writeErr.Error()))
		}
		return
	}

	var req upsertExerciseRequest
	if ok := h.decodeAndValidate(w, r, &req); !ok {
		return
	}

	exercise := &models.Exercise{
		ID:      id,
		Text:    req.Text,
		LevelID: req.LevelID,
	}

	if err := h.exercisesService.Update(r.Context(), exercise); err != nil {
		h.log.Warn("update exercise failed", slog.String("exercise_id", id.String()), slog.String("error", err.Error()))
		apiErr := utils.MapError(err)
		if writeErr := utils.WriteError(w, apiErr.Status, apiErr.Code, apiErr.Message); writeErr != nil {
			h.log.Error("update exercise: write service error response failed", slog.Int("status", apiErr.Status), slog.String("code", string(apiErr.Code)), slog.String("error", writeErr.Error()))
		}
		return
	}

	if err := utils.WriteJSON(w, http.StatusOK, exerciseSingleResponse{Exercise: toExerciseResponse(exercise)}); err != nil {
		h.log.Error("update exercise: write success response failed", slog.String("exercise_id", id.String()), slog.String("error", err.Error()))
	}
}
