package statistics

import (
	"log/slog"
	"net/http"

	"backend-typing-trainer/internal/domain/models"
	"backend-typing-trainer/internal/utils"
)

// Create godoc
// @Summary Создание записи статистики
// @Tags Statistics
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body createStatisticRequest true "Statistic payload"
// @Success 201 {object} statisticSingleResponse "Статистика создана"
// @Failure 400 {object} utils.ErrorResponse "Неверный запрос"
// @Failure 401 {object} utils.ErrorResponse "Не авторизован"
// @Failure 500 {object} utils.InternalErrorResponse "Внутренняя ошибка сервера"
// @Router /statistics [post]
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	h.log.Debug("create statistic started")

	if h.statisticsService == nil {
		h.log.Error("create statistic failed: service is nil")
		if err := utils.WriteError(w, http.StatusInternalServerError, utils.ErrorCodeInternalError, "internal server error"); err != nil {
			h.log.Error("create statistic: write 500 response failed", slog.String("error", err.Error()))
		}
		return
	}

	var req createStatisticRequest
	if ok := h.decodeAndValidate(w, r, &req); !ok {
		return
	}

	statistic := &models.Statistic{
		UserID:          req.UserID,
		LevelID:         req.LevelID,
		ExerciseID:      req.ExerciseID,
		MistakesPercent: req.MistakesPercent,
		ExecutionTime:   req.ExecutionTime,
		Speed:           req.Speed,
	}

	if err := h.statisticsService.Create(r.Context(), statistic); err != nil {
		h.log.Warn("create statistic failed", slog.String("error", err.Error()))
		apiErr := utils.MapError(err)
		if writeErr := utils.WriteError(w, apiErr.Status, apiErr.Code, apiErr.Message); writeErr != nil {
			h.log.Error("create statistic: write service error response failed", slog.Int("status", apiErr.Status), slog.String("code", string(apiErr.Code)), slog.String("error", writeErr.Error()))
		}
		return
	}

	if err := utils.WriteJSON(w, http.StatusCreated, statisticSingleResponse{Statistic: toStatisticResponse(statistic)}); err != nil {
		h.log.Error("create statistic: write success response failed", slog.String("error", err.Error()))
	}
}
