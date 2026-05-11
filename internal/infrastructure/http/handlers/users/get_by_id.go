package users

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"backend-typing-trainer/internal/utils"
)

// GetByID godoc
// @Summary Получение пользователя по ID
// @Tags Users
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID" format(uuid)
// @Success 200 {object} userSingleResponse "Пользователь"
// @Failure 400 {object} utils.ErrorResponse "Неверный запрос"
// @Failure 401 {object} utils.ErrorResponse "Не авторизован"
// @Failure 403 {object} utils.ErrorResponse "Нет доступа"
// @Failure 404 {object} utils.ErrorResponse "Пользователь не найден"
// @Failure 500 {object} utils.InternalErrorResponse "Внутренняя ошибка сервера"
// @Router /users/{id} [get]
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	h.log.Debug("get user by id started")

	if h.usersService == nil {
		h.log.Error("get user by id failed: service is nil")
		if err := utils.WriteError(w, http.StatusInternalServerError, utils.ErrorCodeInternalError, "internal server error"); err != nil {
			h.log.Error("get user by id: write 500 response failed", slog.String("error", err.Error()))
		}
		return
	}

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		h.log.Warn("get user by id failed: invalid id", slog.String("error", err.Error()))
		apiErr := utils.MapError(utils.ErrInvalidRequest)
		if writeErr := utils.WriteError(w, apiErr.Status, apiErr.Code, apiErr.Message); writeErr != nil {
			h.log.Error("get user by id: write invalid id error response failed", slog.Int("status", apiErr.Status), slog.String("code", string(apiErr.Code)), slog.String("error", writeErr.Error()))
		}
		return
	}

	user, err := h.usersService.GetByID(r.Context(), id)
	if err != nil {
		h.log.Warn("get user by id failed", slog.String("user_id", id.String()), slog.String("error", err.Error()))
		apiErr := utils.MapError(err)
		if writeErr := utils.WriteError(w, apiErr.Status, apiErr.Code, apiErr.Message); writeErr != nil {
			h.log.Error("get user by id: write service error response failed", slog.Int("status", apiErr.Status), slog.String("code", string(apiErr.Code)), slog.String("error", writeErr.Error()))
		}
		return
	}

	if err := utils.WriteJSON(w, http.StatusOK, userSingleResponse{User: fromModel(user)}); err != nil {
		h.log.Error("get user by id: write success response failed", slog.String("user_id", id.String()), slog.String("error", err.Error()))
	}
}
