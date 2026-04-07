package keyboardzones

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	"backend-typing-trainer/internal/utils"
)

// GetByName godoc
// @Summary Получение зоны клавиатуры по имени
// @Tags KeyboardZones
// @Produce json
// @Security BearerAuth
// @Param name path string true "Keyboard zone name"
// @Success 200 {object} keyboardZoneSingleResponse "Зона клавиатуры"
// @Failure 400 {object} utils.ErrorResponse "Неверный запрос"
// @Failure 401 {object} utils.ErrorResponse "Не авторизован"
// @Failure 404 {object} utils.ErrorResponse "Зона не найдена"
// @Failure 500 {object} utils.InternalErrorResponse "Внутренняя ошибка сервера"
// @Router /keyboard-zones/by-name/{name} [get]
func (h *Handler) GetByName(w http.ResponseWriter, r *http.Request) {
	h.log.Debug("get keyboard zone by name started")

	if h.keyboardZonesService == nil {
		h.log.Error("get keyboard zone by name failed: service is nil")
		if err := utils.WriteError(w, http.StatusInternalServerError, utils.ErrorCodeInternalError, "internal server error"); err != nil {
			h.log.Error("get keyboard zone by name: write 500 response failed", slog.String("error", err.Error()))
		}
		return
	}

	name := strings.TrimSpace(chi.URLParam(r, "name"))
	if name == "" {
		h.log.Warn("get keyboard zone by name failed: empty name")
		apiErr := utils.MapError(utils.ErrInvalidRequest)
		if writeErr := utils.WriteError(w, apiErr.Status, apiErr.Code, apiErr.Message); writeErr != nil {
			h.log.Error("get keyboard zone by name: write invalid name error response failed", slog.Int("status", apiErr.Status), slog.String("code", string(apiErr.Code)), slog.String("error", writeErr.Error()))
		}
		return
	}

	zone, err := h.keyboardZonesService.GetByName(r.Context(), name)
	if err != nil {
		h.log.Warn("get keyboard zone by name failed", slog.String("name", name), slog.String("error", err.Error()))
		apiErr := utils.MapError(err)
		if writeErr := utils.WriteError(w, apiErr.Status, apiErr.Code, apiErr.Message); writeErr != nil {
			h.log.Error("get keyboard zone by name: write service error response failed", slog.Int("status", apiErr.Status), slog.String("code", string(apiErr.Code)), slog.String("error", writeErr.Error()))
		}
		return
	}

	if err := utils.WriteJSON(w, http.StatusOK, keyboardZoneSingleResponse{KeyboardZone: toKeyboardZoneResponse(zone)}); err != nil {
		h.log.Error("get keyboard zone by name: write success response failed", slog.String("name", name), slog.String("error", err.Error()))
	}
}
