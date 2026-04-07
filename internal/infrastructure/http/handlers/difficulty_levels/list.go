package difficultylevels

import (
	"log/slog"
	"net/http"
	"strconv"

	"backend-typing-trainer/internal/utils"
)

const (
	defaultListLimit  = 20
	defaultListOffset = 0
)

// List godoc
// @Summary Список уровней сложности
// @Tags DifficultyLevels
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Лимит" default(20) minimum(1)
// @Param offset query int false "Смещение" default(0) minimum(0)
// @Success 200 {object} difficultyLevelListResponse "Список уровней сложности"
// @Failure 400 {object} utils.ErrorResponse "Неверный запрос"
// @Failure 401 {object} utils.ErrorResponse "Не авторизован"
// @Failure 500 {object} utils.InternalErrorResponse "Внутренняя ошибка сервера"
// @Router /difficulty-levels [get]
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	h.log.Debug("list difficulty levels started")

	if h.difficultyLevelsService == nil {
		h.log.Error("list difficulty levels failed: service is nil")
		if err := utils.WriteError(w, http.StatusInternalServerError, utils.ErrorCodeInternalError, "internal server error"); err != nil {
			h.log.Error("list difficulty levels: write 500 response failed", slog.String("error", err.Error()))
		}
		return
	}

	limit, offset, ok := parsePagination(r)
	if !ok {
		h.log.Warn("list difficulty levels failed: invalid pagination")
		apiErr := utils.MapError(utils.ErrInvalidRequest)
		if writeErr := utils.WriteError(w, apiErr.Status, apiErr.Code, apiErr.Message); writeErr != nil {
			h.log.Error("list difficulty levels: write invalid pagination error response failed", slog.Int("status", apiErr.Status), slog.String("code", string(apiErr.Code)), slog.String("error", writeErr.Error()))
		}
		return
	}

	levels, err := h.difficultyLevelsService.List(r.Context(), limit, offset)
	if err != nil {
		h.log.Warn("list difficulty levels failed", slog.String("error", err.Error()))
		apiErr := utils.MapError(err)
		if writeErr := utils.WriteError(w, apiErr.Status, apiErr.Code, apiErr.Message); writeErr != nil {
			h.log.Error("list difficulty levels: write service error response failed", slog.Int("status", apiErr.Status), slog.String("code", string(apiErr.Code)), slog.String("error", writeErr.Error()))
		}
		return
	}

	items := make([]difficultyLevelResponse, 0, len(levels))
	for _, level := range levels {
		items = append(items, toDifficultyLevelResponse(level))
	}

	if err := utils.WriteJSON(w, http.StatusOK, difficultyLevelListResponse{DifficultyLevels: items}); err != nil {
		h.log.Error("list difficulty levels: write success response failed", slog.String("error", err.Error()))
	}
}

func parsePagination(r *http.Request) (int, int, bool) {
	limit := defaultListLimit
	offset := defaultListOffset

	if rawLimit := r.URL.Query().Get("limit"); rawLimit != "" {
		parsedLimit, err := strconv.Atoi(rawLimit)
		if err != nil {
			return 0, 0, false
		}
		limit = parsedLimit
	}

	if rawOffset := r.URL.Query().Get("offset"); rawOffset != "" {
		parsedOffset, err := strconv.Atoi(rawOffset)
		if err != nil {
			return 0, 0, false
		}
		offset = parsedOffset
	}

	if limit <= 0 || offset < 0 {
		return 0, 0, false
	}

	return limit, offset, true
}
