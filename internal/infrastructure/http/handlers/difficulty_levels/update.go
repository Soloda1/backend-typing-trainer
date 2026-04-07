package difficultylevels

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"backend-typing-trainer/internal/domain/models"
	"backend-typing-trainer/internal/utils"
)

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	h.log.Debug("update difficulty level started")

	if h.difficultyLevelsService == nil {
		h.log.Error("update difficulty level failed: service is nil")
		if err := utils.WriteError(w, http.StatusInternalServerError, utils.ErrorCodeInternalError, "internal server error"); err != nil {
			h.log.Error("update difficulty level: write 500 response failed", slog.String("error", err.Error()))
		}
		return
	}

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		h.log.Warn("update difficulty level failed: invalid id", slog.String("error", err.Error()))
		apiErr := utils.MapError(utils.ErrInvalidRequest)
		if writeErr := utils.WriteError(w, apiErr.Status, apiErr.Code, apiErr.Message); writeErr != nil {
			h.log.Error("update difficulty level: write invalid id error response failed", slog.Int("status", apiErr.Status), slog.String("code", string(apiErr.Code)), slog.String("error", writeErr.Error()))
		}
		return
	}

	var req upsertDifficultyLevelRequest
	if ok := h.decodeAndValidate(w, r, &req); !ok {
		return
	}

	level := &models.DifficultyLevel{
		ID:                id,
		AllowedMistakes:   req.AllowedMistakes,
		KeyPressTime:      req.KeyPressTime,
		MinExerciseLength: req.MinExerciseLength,
		MaxExerciseLength: req.MaxExerciseLength,
		KeyboardZoneIDs:   req.KeyboardZoneIDs,
	}

	if err := h.difficultyLevelsService.Update(r.Context(), level); err != nil {
		h.log.Warn("update difficulty level failed", slog.String("level_id", id.String()), slog.String("error", err.Error()))
		apiErr := utils.MapError(err)
		if writeErr := utils.WriteError(w, apiErr.Status, apiErr.Code, apiErr.Message); writeErr != nil {
			h.log.Error("update difficulty level: write service error response failed", slog.Int("status", apiErr.Status), slog.String("code", string(apiErr.Code)), slog.String("error", writeErr.Error()))
		}
		return
	}

	if err := utils.WriteJSON(w, http.StatusOK, difficultyLevelSingleResponse{DifficultyLevel: toDifficultyLevelResponse(level)}); err != nil {
		h.log.Error("update difficulty level: write success response failed", slog.String("level_id", id.String()), slog.String("error", err.Error()))
	}
}
