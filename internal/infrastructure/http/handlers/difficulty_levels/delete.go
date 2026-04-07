package difficultylevels

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"backend-typing-trainer/internal/utils"
)

// Delete godoc
// @Summary Удаление уровня сложности
// @Description Доступно только администратору.
// @Tags DifficultyLevels
// @Produce json
// @Security BearerAuth
// @Param id path string true "Difficulty level ID" format(uuid)
// @Success 204 "Уровень сложности удален"
// @Failure 400 {object} utils.ErrorResponse "Неверный запрос"
// @Failure 401 {object} utils.ErrorResponse "Не авторизован"
// @Failure 403 {object} utils.ErrorResponse "Недостаточно прав"
// @Failure 404 {object} utils.ErrorResponse "Уровень не найден"
// @Failure 500 {object} utils.InternalErrorResponse "Внутренняя ошибка сервера"
// @Router /difficulty-levels/{id} [delete]
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	h.log.Debug("delete difficulty level started")

	if h.difficultyLevelsService == nil {
		h.log.Error("delete difficulty level failed: service is nil")
		if err := utils.WriteError(w, http.StatusInternalServerError, utils.ErrorCodeInternalError, "internal server error"); err != nil {
			h.log.Error("delete difficulty level: write 500 response failed", slog.String("error", err.Error()))
		}
		return
	}

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		h.log.Warn("delete difficulty level failed: invalid id", slog.String("error", err.Error()))
		apiErr := utils.MapError(utils.ErrInvalidRequest)
		if writeErr := utils.WriteError(w, apiErr.Status, apiErr.Code, apiErr.Message); writeErr != nil {
			h.log.Error("delete difficulty level: write invalid id error response failed", slog.Int("status", apiErr.Status), slog.String("code", string(apiErr.Code)), slog.String("error", writeErr.Error()))
		}
		return
	}

	if err := h.difficultyLevelsService.Delete(r.Context(), id); err != nil {
		h.log.Warn("delete difficulty level failed", slog.String("level_id", id.String()), slog.String("error", err.Error()))
		apiErr := utils.MapError(err)
		if writeErr := utils.WriteError(w, apiErr.Status, apiErr.Code, apiErr.Message); writeErr != nil {
			h.log.Error("delete difficulty level: write service error response failed", slog.Int("status", apiErr.Status), slog.String("code", string(apiErr.Code)), slog.String("error", writeErr.Error()))
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
