package exercises

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"backend-typing-trainer/internal/utils"
)

// Delete godoc
// @Summary Удаление упражнения
// @Description Доступно только администратору.
// @Tags Exercises
// @Produce json
// @Security BearerAuth
// @Param id path string true "Exercise ID" format(uuid)
// @Success 204 "Упражнение удалено"
// @Failure 400 {object} utils.ErrorResponse "Неверный запрос"
// @Failure 401 {object} utils.ErrorResponse "Не авторизован"
// @Failure 403 {object} utils.ErrorResponse "Недостаточно прав"
// @Failure 404 {object} utils.ErrorResponse "Упражнение не найдено"
// @Failure 500 {object} utils.InternalErrorResponse "Внутренняя ошибка сервера"
// @Router /exercises/{id} [delete]
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	h.log.Debug("delete exercise started")

	if h.exercisesService == nil {
		h.log.Error("delete exercise failed: service is nil")
		if err := utils.WriteError(w, http.StatusInternalServerError, utils.ErrorCodeInternalError, "internal server error"); err != nil {
			h.log.Error("delete exercise: write 500 response failed", slog.String("error", err.Error()))
		}
		return
	}

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		h.log.Warn("delete exercise failed: invalid id", slog.String("error", err.Error()))
		apiErr := utils.MapError(utils.ErrInvalidRequest)
		if writeErr := utils.WriteError(w, apiErr.Status, apiErr.Code, apiErr.Message); writeErr != nil {
			h.log.Error("delete exercise: write invalid id error response failed", slog.Int("status", apiErr.Status), slog.String("code", string(apiErr.Code)), slog.String("error", writeErr.Error()))
		}
		return
	}

	if err := h.exercisesService.Delete(r.Context(), id); err != nil {
		h.log.Warn("delete exercise failed", slog.String("exercise_id", id.String()), slog.String("error", err.Error()))
		apiErr := utils.MapError(err)
		if writeErr := utils.WriteError(w, apiErr.Status, apiErr.Code, apiErr.Message); writeErr != nil {
			h.log.Error("delete exercise: write service error response failed", slog.Int("status", apiErr.Status), slog.String("code", string(apiErr.Code)), slog.String("error", writeErr.Error()))
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
