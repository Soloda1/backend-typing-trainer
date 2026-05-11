package users

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	"backend-typing-trainer/internal/utils"
)

// GetByLogin godoc
// @Summary Получение пользователя по логину
// @Tags Users
// @Produce json
// @Security BearerAuth
// @Param login path string true "User Login"
// @Success 200 {object} userSingleResponse "Пользователь"
// @Failure 400 {object} utils.ErrorResponse "Неверный запрос"
// @Failure 401 {object} utils.ErrorResponse "Не авторизован"
// @Failure 403 {object} utils.ErrorResponse "Нет доступа"
// @Failure 404 {object} utils.ErrorResponse "Пользователь не найден"
// @Failure 500 {object} utils.InternalErrorResponse "Внутренняя ошибка сервера"
// @Router /users/login/{login} [get]
func (h *Handler) GetByLogin(w http.ResponseWriter, r *http.Request) {
	h.log.Debug("get user by login started")

	if h.usersService == nil {
		h.log.Error("get user by login failed: service is nil")
		if err := utils.WriteError(w, http.StatusInternalServerError, utils.ErrorCodeInternalError, "internal server error"); err != nil {
			h.log.Error("get user by login: write 500 response failed", slog.String("error", err.Error()))
		}
		return
	}

	login := strings.TrimSpace(chi.URLParam(r, "login"))
	if login == "" {
		h.log.Warn("get user by login failed: empty login")
		apiErr := utils.MapError(utils.ErrInvalidRequest)
		if writeErr := utils.WriteError(w, apiErr.Status, apiErr.Code, apiErr.Message); writeErr != nil {
			h.log.Error("get user by login: write invalid login error response failed", slog.Int("status", apiErr.Status), slog.String("code", string(apiErr.Code)), slog.String("error", writeErr.Error()))
		}
		return
	}

	user, err := h.usersService.GetByLogin(r.Context(), login)
	if err != nil {
		h.log.Warn("get user by login failed", slog.String("login", login), slog.String("error", err.Error()))
		apiErr := utils.MapError(err)
		if writeErr := utils.WriteError(w, apiErr.Status, apiErr.Code, apiErr.Message); writeErr != nil {
			h.log.Error("get user by login: write service error response failed", slog.Int("status", apiErr.Status), slog.String("code", string(apiErr.Code)), slog.String("error", writeErr.Error()))
		}
		return
	}

	if err := utils.WriteJSON(w, http.StatusOK, userSingleResponse{User: fromModel(user)}); err != nil {
		h.log.Error("get user by login: write success response failed", slog.String("login", login), slog.String("error", err.Error()))
	}
}
