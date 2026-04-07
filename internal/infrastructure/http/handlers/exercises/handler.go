package exercises

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	input "backend-typing-trainer/internal/domain/ports/input"
	ports "backend-typing-trainer/internal/domain/ports/output/logger"
	"backend-typing-trainer/internal/utils"
)

type Handler struct {
	log              ports.Logger
	exercisesService input.Exercises
}

func NewHandler(log ports.Logger, exercisesService input.Exercises) *Handler {
	log = log.With("handler", "exercises")

	return &Handler{
		log:              log,
		exercisesService: exercisesService,
	}
}

func (h *Handler) decodeAndValidate(w http.ResponseWriter, r *http.Request, dst any) bool {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(dst); err != nil {
		h.log.Warn("invalid request body", slog.String("error", err.Error()))
		apiErr := utils.MapError(utils.ErrInvalidRequest)
		if writeErr := utils.WriteError(w, apiErr.Status, apiErr.Code, apiErr.Message); writeErr != nil {
			h.log.Error("write invalid body error response failed", slog.Int("status", apiErr.Status), slog.String("code", string(apiErr.Code)), slog.String("error", writeErr.Error()))
		}
		return false
	}

	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		if err == nil {
			h.log.Warn("invalid request body: extra json payload")
		} else {
			h.log.Warn("invalid request body: extra json payload", slog.String("error", err.Error()))
		}
		apiErr := utils.MapError(utils.ErrInvalidRequest)
		if writeErr := utils.WriteError(w, apiErr.Status, apiErr.Code, apiErr.Message); writeErr != nil {
			h.log.Error("write extra payload error response failed", slog.Int("status", apiErr.Status), slog.String("code", string(apiErr.Code)), slog.String("error", writeErr.Error()))
		}
		return false
	}

	if err := utils.Validate(dst); err != nil {
		h.log.Warn("request validation failed", slog.String("error", err.Error()))
		apiErr := utils.MapError(utils.ErrInvalidRequest)
		if writeErr := utils.WriteError(w, apiErr.Status, apiErr.Code, apiErr.Message); writeErr != nil {
			h.log.Error("write validation error response failed", slog.Int("status", apiErr.Status), slog.String("code", string(apiErr.Code)), slog.String("error", writeErr.Error()))
		}
		return false
	}

	return true
}
