package statistics

import (
	"backend-typing-trainer/internal/domain/models"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"backend-typing-trainer/internal/utils"
)

const (
	defaultListLimit  = 20
	defaultListOffset = 0
)

// ListByUserID godoc
// @Summary Список статистики пользователя
// @Tags Statistics
// @Produce json
// @Security BearerAuth
// @Param user_id path string true "User ID" format(uuid)
// @Param limit query int false "Лимит" default(20) minimum(1)
// @Param offset query int false "Смещение" default(0) minimum(0)
// @Success 200 {object} statisticListResponse "Список статистики"
// @Failure 400 {object} utils.ErrorResponse "Неверный запрос"
// @Failure 401 {object} utils.ErrorResponse "Не авторизован"
// @Failure 403 {object} utils.ErrorResponse "Недостаточно прав"
// @Failure 500 {object} utils.InternalErrorResponse "Внутренняя ошибка сервера"
// @Router /statistics/users/{user_id} [get]
func (h *Handler) ListByUserID(w http.ResponseWriter, r *http.Request) {
	h.log.Debug("list statistics by user started")

	if h.statisticsService == nil {
		h.log.Error("list statistics by user failed: service is nil")
		if err := utils.WriteError(w, http.StatusInternalServerError, utils.ErrorCodeInternalError, "internal server error"); err != nil {
			h.log.Error("list statistics by user: write 500 response failed", slog.String("error", err.Error()))
		}
		return
	}

	userID, err := uuid.Parse(chi.URLParam(r, "user_id"))
	if err != nil {
		h.log.Warn("list statistics by user failed: invalid user id", slog.String("error", err.Error()))
		apiErr := utils.MapError(utils.ErrInvalidRequest)
		if writeErr := utils.WriteError(w, apiErr.Status, apiErr.Code, apiErr.Message); writeErr != nil {
			h.log.Error("list statistics by user: write invalid id error response failed", slog.Int("status", apiErr.Status), slog.String("code", string(apiErr.Code)), slog.String("error", writeErr.Error()))
		}
		return
	}

	limit, offset, ok := parsePagination(r)
	if !ok {
		h.log.Warn("list statistics by user failed: invalid pagination")
		apiErr := utils.MapError(utils.ErrInvalidRequest)
		if writeErr := utils.WriteError(w, apiErr.Status, apiErr.Code, apiErr.Message); writeErr != nil {
			h.log.Error("list statistics by user: write invalid pagination error response failed", slog.Int("status", apiErr.Status), slog.String("code", string(apiErr.Code)), slog.String("error", writeErr.Error()))
		}
		return
	}

	statistics, err := h.statisticsService.ListByUserID(r.Context(), userID, limit, offset)
	if err != nil {
		h.log.Warn("list statistics by user failed", slog.String("user_id", userID.String()), slog.String("error", err.Error()))
		apiErr := utils.MapError(err)
		if writeErr := utils.WriteError(w, apiErr.Status, apiErr.Code, apiErr.Message); writeErr != nil {
			h.log.Error("list statistics by user: write service error response failed", slog.Int("status", apiErr.Status), slog.String("code", string(apiErr.Code)), slog.String("error", writeErr.Error()))
		}
		return
	}

	h.writeStatisticList(w, statistics)
}

// ListByLevelID godoc
// @Summary Список статистики по уровню сложности
// @Tags Statistics
// @Produce json
// @Security BearerAuth
// @Param level_id path string true "Level ID" format(uuid)
// @Param limit query int false "Лимит" default(20) minimum(1)
// @Param offset query int false "Смещение" default(0) minimum(0)
// @Success 200 {object} statisticListResponse "Список статистики"
// @Failure 400 {object} utils.ErrorResponse "Неверный запрос"
// @Failure 401 {object} utils.ErrorResponse "Не авторизован"
// @Failure 403 {object} utils.ErrorResponse "Недостаточно прав"
// @Failure 500 {object} utils.InternalErrorResponse "Внутренняя ошибка сервера"
// @Router /statistics/levels/{level_id} [get]
func (h *Handler) ListByLevelID(w http.ResponseWriter, r *http.Request) {
	h.log.Debug("list statistics by level started")

	if h.statisticsService == nil {
		h.log.Error("list statistics by level failed: service is nil")
		if err := utils.WriteError(w, http.StatusInternalServerError, utils.ErrorCodeInternalError, "internal server error"); err != nil {
			h.log.Error("list statistics by level: write 500 response failed", slog.String("error", err.Error()))
		}
		return
	}

	levelID, err := uuid.Parse(chi.URLParam(r, "level_id"))
	if err != nil {
		h.log.Warn("list statistics by level failed: invalid level id", slog.String("error", err.Error()))
		apiErr := utils.MapError(utils.ErrInvalidRequest)
		if writeErr := utils.WriteError(w, apiErr.Status, apiErr.Code, apiErr.Message); writeErr != nil {
			h.log.Error("list statistics by level: write invalid id error response failed", slog.Int("status", apiErr.Status), slog.String("code", string(apiErr.Code)), slog.String("error", writeErr.Error()))
		}
		return
	}

	limit, offset, ok := parsePagination(r)
	if !ok {
		h.log.Warn("list statistics by level failed: invalid pagination")
		apiErr := utils.MapError(utils.ErrInvalidRequest)
		if writeErr := utils.WriteError(w, apiErr.Status, apiErr.Code, apiErr.Message); writeErr != nil {
			h.log.Error("list statistics by level: write invalid pagination error response failed", slog.Int("status", apiErr.Status), slog.String("code", string(apiErr.Code)), slog.String("error", writeErr.Error()))
		}
		return
	}

	statistics, err := h.statisticsService.ListByLevelID(r.Context(), levelID, limit, offset)
	if err != nil {
		h.log.Warn("list statistics by level failed", slog.String("level_id", levelID.String()), slog.String("error", err.Error()))
		apiErr := utils.MapError(err)
		if writeErr := utils.WriteError(w, apiErr.Status, apiErr.Code, apiErr.Message); writeErr != nil {
			h.log.Error("list statistics by level: write service error response failed", slog.Int("status", apiErr.Status), slog.String("code", string(apiErr.Code)), slog.String("error", writeErr.Error()))
		}
		return
	}

	h.writeStatisticList(w, statistics)
}

// ListByExerciseID godoc
// @Summary Список статистики по упражнению
// @Tags Statistics
// @Produce json
// @Security BearerAuth
// @Param exercise_id path string true "Exercise ID" format(uuid)
// @Param limit query int false "Лимит" default(20) minimum(1)
// @Param offset query int false "Смещение" default(0) minimum(0)
// @Success 200 {object} statisticListResponse "Список статистики"
// @Failure 400 {object} utils.ErrorResponse "Неверный запрос"
// @Failure 401 {object} utils.ErrorResponse "Не авторизован"
// @Failure 403 {object} utils.ErrorResponse "Недостаточно прав"
// @Failure 500 {object} utils.InternalErrorResponse "Внутренняя ошибка сервера"
// @Router /statistics/exercises/{exercise_id} [get]
func (h *Handler) ListByExerciseID(w http.ResponseWriter, r *http.Request) {
	h.log.Debug("list statistics by exercise started")

	if h.statisticsService == nil {
		h.log.Error("list statistics by exercise failed: service is nil")
		if err := utils.WriteError(w, http.StatusInternalServerError, utils.ErrorCodeInternalError, "internal server error"); err != nil {
			h.log.Error("list statistics by exercise: write 500 response failed", slog.String("error", err.Error()))
		}
		return
	}

	exerciseID, err := uuid.Parse(chi.URLParam(r, "exercise_id"))
	if err != nil {
		h.log.Warn("list statistics by exercise failed: invalid exercise id", slog.String("error", err.Error()))
		apiErr := utils.MapError(utils.ErrInvalidRequest)
		if writeErr := utils.WriteError(w, apiErr.Status, apiErr.Code, apiErr.Message); writeErr != nil {
			h.log.Error("list statistics by exercise: write invalid id error response failed", slog.Int("status", apiErr.Status), slog.String("code", string(apiErr.Code)), slog.String("error", writeErr.Error()))
		}
		return
	}

	limit, offset, ok := parsePagination(r)
	if !ok {
		h.log.Warn("list statistics by exercise failed: invalid pagination")
		apiErr := utils.MapError(utils.ErrInvalidRequest)
		if writeErr := utils.WriteError(w, apiErr.Status, apiErr.Code, apiErr.Message); writeErr != nil {
			h.log.Error("list statistics by exercise: write invalid pagination error response failed", slog.Int("status", apiErr.Status), slog.String("code", string(apiErr.Code)), slog.String("error", writeErr.Error()))
		}
		return
	}

	statistics, err := h.statisticsService.ListByExerciseID(r.Context(), exerciseID, limit, offset)
	if err != nil {
		h.log.Warn("list statistics by exercise failed", slog.String("exercise_id", exerciseID.String()), slog.String("error", err.Error()))
		apiErr := utils.MapError(err)
		if writeErr := utils.WriteError(w, apiErr.Status, apiErr.Code, apiErr.Message); writeErr != nil {
			h.log.Error("list statistics by exercise: write service error response failed", slog.Int("status", apiErr.Status), slog.String("code", string(apiErr.Code)), slog.String("error", writeErr.Error()))
		}
		return
	}

	h.writeStatisticList(w, statistics)
}

func (h *Handler) writeStatisticList(w http.ResponseWriter, statistics []*models.Statistic) {

	items := make([]statisticResponse, 0, len(statistics))
	for _, statistic := range statistics {
		items = append(items, toStatisticResponse(statistic))
	}

	if err := utils.WriteJSON(w, http.StatusOK, statisticListResponse{Statistics: items}); err != nil {
		h.log.Error("list statistics: write success response failed", slog.String("error", err.Error()))
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
