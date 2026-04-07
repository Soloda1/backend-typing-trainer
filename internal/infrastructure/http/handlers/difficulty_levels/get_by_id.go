package difficultylevels

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"backend-typing-trainer/internal/utils"
)

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	h.log.Debug("get difficulty level by id started")

	if h.difficultyLevelsService == nil {
		h.log.Error("get difficulty level by id failed: service is nil")
		if err := utils.WriteError(w, http.StatusInternalServerError, utils.ErrorCodeInternalError, "internal server error"); err != nil {
			h.log.Error("get difficulty level by id: write 500 response failed", slog.String("error", err.Error()))
		}
		return
	}

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		h.log.Warn("get difficulty level by id failed: invalid id", slog.String("error", err.Error()))
		apiErr := utils.MapError(utils.ErrInvalidRequest)
		if writeErr := utils.WriteError(w, apiErr.Status, apiErr.Code, apiErr.Message); writeErr != nil {
			h.log.Error("get difficulty level by id: write invalid id error response failed", slog.Int("status", apiErr.Status), slog.String("code", string(apiErr.Code)), slog.String("error", writeErr.Error()))
		}
		return
	}

	level, err := h.difficultyLevelsService.GetByID(r.Context(), id)
	if err != nil {
		h.log.Warn("get difficulty level by id failed", slog.String("level_id", id.String()), slog.String("error", err.Error()))
		apiErr := utils.MapError(err)
		if writeErr := utils.WriteError(w, apiErr.Status, apiErr.Code, apiErr.Message); writeErr != nil {
			h.log.Error("get difficulty level by id: write service error response failed", slog.Int("status", apiErr.Status), slog.String("code", string(apiErr.Code)), slog.String("error", writeErr.Error()))
		}
		return
	}

	if err := utils.WriteJSON(w, http.StatusOK, difficultyLevelSingleResponse{DifficultyLevel: toDifficultyLevelResponse(level)}); err != nil {
		h.log.Error("get difficulty level by id: write success response failed", slog.String("level_id", id.String()), slog.String("error", err.Error()))
	}
}
