package statistics

import (
	"backend-typing-trainer/internal/domain/models"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/google/uuid"

	"backend-typing-trainer/internal/utils"
)

const (
	defaultListLimit  = 20
	defaultListOffset = 0
)

// List godoc
// @Summary Список статистики по фильтру
// @Description Нужно передать ровно один фильтр: user_id, level_id или exercise_id.
// @Tags Statistics
// @Produce json
// @Security BearerAuth
// @Param user_id query string false "User ID" format(uuid)
// @Param level_id query string false "Level ID" format(uuid)
// @Param exercise_id query string false "Exercise ID" format(uuid)
// @Param limit query int false "Лимит" default(20) minimum(1)
// @Param offset query int false "Смещение" default(0) minimum(0)
// @Success 200 {object} statisticListResponse "Список статистики"
// @Failure 400 {object} utils.ErrorResponse "Неверный запрос"
// @Failure 401 {object} utils.ErrorResponse "Не авторизован"
// @Failure 500 {object} utils.InternalErrorResponse "Внутренняя ошибка сервера"
// @Router /statistics [get]
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	h.log.Debug("list statistics started")

	if h.statisticsService == nil {
		h.log.Error("list statistics failed: service is nil")
		if err := utils.WriteError(w, http.StatusInternalServerError, utils.ErrorCodeInternalError, "internal server error"); err != nil {
			h.log.Error("list statistics: write 500 response failed", slog.String("error", err.Error()))
		}
		return
	}

	filterID, filterType, limit, offset, ok := parseListFilters(r)
	if !ok {
		h.log.Warn("list statistics failed: invalid filters")
		apiErr := utils.MapError(utils.ErrInvalidRequest)
		if writeErr := utils.WriteError(w, apiErr.Status, apiErr.Code, apiErr.Message); writeErr != nil {
			h.log.Error("list statistics: write invalid filters error response failed", slog.Int("status", apiErr.Status), slog.String("code", string(apiErr.Code)), slog.String("error", writeErr.Error()))
		}
		return
	}

	var (
		statistics []*models.Statistic
		err        error
	)

	switch filterType {
	case "user_id":
		statistics, err = h.statisticsService.ListByUserID(r.Context(), filterID, limit, offset)
	case "level_id":
		statistics, err = h.statisticsService.ListByLevelID(r.Context(), filterID, limit, offset)
	case "exercise_id":
		statistics, err = h.statisticsService.ListByExerciseID(r.Context(), filterID, limit, offset)
	default:
		err = utils.ErrInvalidRequest
	}

	if err != nil {
		h.log.Warn("list statistics failed", slog.String("filter", filterType), slog.String("error", err.Error()))
		apiErr := utils.MapError(err)
		if writeErr := utils.WriteError(w, apiErr.Status, apiErr.Code, apiErr.Message); writeErr != nil {
			h.log.Error("list statistics: write service error response failed", slog.Int("status", apiErr.Status), slog.String("code", string(apiErr.Code)), slog.String("error", writeErr.Error()))
		}
		return
	}

	items := make([]statisticResponse, 0, len(statistics))
	for _, statistic := range statistics {
		items = append(items, toStatisticResponse(statistic))
	}

	if err := utils.WriteJSON(w, http.StatusOK, statisticListResponse{Statistics: items}); err != nil {
		h.log.Error("list statistics: write success response failed", slog.String("error", err.Error()))
	}
}

func parseListFilters(r *http.Request) (uuid.UUID, string, int, int, bool) {
	limit := defaultListLimit
	offset := defaultListOffset

	if rawLimit := r.URL.Query().Get("limit"); rawLimit != "" {
		parsedLimit, err := strconv.Atoi(rawLimit)
		if err != nil {
			return uuid.Nil, "", 0, 0, false
		}
		limit = parsedLimit
	}
	if rawOffset := r.URL.Query().Get("offset"); rawOffset != "" {
		parsedOffset, err := strconv.Atoi(rawOffset)
		if err != nil {
			return uuid.Nil, "", 0, 0, false
		}
		offset = parsedOffset
	}
	if limit <= 0 || offset < 0 {
		return uuid.Nil, "", 0, 0, false
	}

	filters := []struct {
		key string
		val string
	}{
		{key: "user_id", val: r.URL.Query().Get("user_id")},
		{key: "level_id", val: r.URL.Query().Get("level_id")},
		{key: "exercise_id", val: r.URL.Query().Get("exercise_id")},
	}

	foundKey := ""
	foundValue := ""
	for _, f := range filters {
		if f.val == "" {
			continue
		}
		if foundKey != "" {
			return uuid.Nil, "", 0, 0, false
		}
		foundKey = f.key
		foundValue = f.val
	}
	if foundKey == "" {
		return uuid.Nil, "", 0, 0, false
	}

	id, err := uuid.Parse(foundValue)
	if err != nil {
		return uuid.Nil, "", 0, 0, false
	}

	return id, foundKey, limit, offset, true
}
