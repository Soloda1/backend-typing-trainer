package exercises

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"backend-typing-trainer/internal/utils"
)

// GetByID godoc
// @Summary Получение упражнения по ID
// @Tags Exercises
// @Produce json
// @Security BearerAuth
// @Param id path string true "Exercise ID" format(uuid)
// @Success 200 {object} exerciseSingleResponse "Упражнение"
// @Failure 400 {object} utils.ErrorResponse "Неверный запрос"
// @Failure 401 {object} utils.ErrorResponse "Не авторизован"
// @Failure 404 {object} utils.ErrorResponse "Упражнение не найдено"
// @Failure 500 {object} utils.InternalErrorResponse "Внутренняя ошибка сервера"
// @Router /exercises/{id} [get]
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	h.log.Debug("get exercise by id started")

	if h.exercisesService == nil {
		h.log.Error("get exercise by id failed: service is nil")
		if err := utils.WriteError(w, http.StatusInternalServerError, utils.ErrorCodeInternalError, "internal server error"); err != nil {
			h.log.Error("get exercise by id: write 500 response failed", slog.String("error", err.Error()))
		}
		return
	}

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		h.log.Warn("get exercise by id failed: invalid id", slog.String("error", err.Error()))
		apiErr := utils.MapError(utils.ErrInvalidRequest)
		if writeErr := utils.WriteError(w, apiErr.Status, apiErr.Code, apiErr.Message); writeErr != nil {
			h.log.Error("get exercise by id: write invalid id error response failed", slog.Int("status", apiErr.Status), slog.String("code", string(apiErr.Code)), slog.String("error", writeErr.Error()))
		}
		return
	}

	exercise, err := h.exercisesService.GetByID(r.Context(), id)
	if err != nil {
		h.log.Warn("get exercise by id failed", slog.String("exercise_id", id.String()), slog.String("error", err.Error()))
		apiErr := utils.MapError(err)
		if writeErr := utils.WriteError(w, apiErr.Status, apiErr.Code, apiErr.Message); writeErr != nil {
			h.log.Error("get exercise by id: write service error response failed", slog.Int("status", apiErr.Status), slog.String("code", string(apiErr.Code)), slog.String("error", writeErr.Error()))
		}
		return
	}

	if err := utils.WriteJSON(w, http.StatusOK, exerciseSingleResponse{Exercise: toExerciseResponse(exercise)}); err != nil {
		h.log.Error("get exercise by id: write success response failed", slog.String("exercise_id", id.String()), slog.String("error", err.Error()))
	}
}
