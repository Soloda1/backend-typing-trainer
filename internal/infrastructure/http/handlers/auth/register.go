package auth

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"

	"backend-typing-trainer/internal/domain/models"
	"backend-typing-trainer/internal/utils"
)

type registerRequest struct {
	Login    string `json:"login" validate:"required" example:"player_one"`
	Password string `json:"password" validate:"required" example:"secret"`
}

type registerUserResponse struct {
	ID        uuid.UUID       `json:"id" example:"22222222-2222-2222-2222-222222222222"`
	Login     string          `json:"login" example:"player_one"`
	Role      models.UserRole `json:"role" example:"user"`
	CreatedAt *time.Time      `json:"createdAt,omitempty" swaggertype:"string" example:"2026-03-22T21:00:00Z"`
}

type registerResponse struct {
	User registerUserResponse `json:"user"`
}

// Register godoc
// @Summary Регистрация пользователя
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body registerRequest true "Register payload"
// @Success 201 {object} registerResponse "Пользователь создан"
// @Failure 400 {object} utils.ErrorResponse "Неверный запрос или логин уже занят"
// @Failure 500 {object} utils.InternalErrorResponse "Внутренняя ошибка сервера"
// @Router /register [post]
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	h.log.Debug("register started")

	if h.authService == nil {
		h.log.Error("register failed: auth service is nil")
		if err := utils.WriteError(w, http.StatusInternalServerError, utils.ErrorCodeInternalError, "internal server error"); err != nil {
			h.log.Error("register: write 500 response failed", slog.String("error", err.Error()))
		}
		return
	}

	var req registerRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		h.log.Warn("register invalid request body", slog.String("error", err.Error()))
		apiErr := utils.MapError(utils.ErrInvalidRequest)
		if err := utils.WriteError(w, apiErr.Status, apiErr.Code, apiErr.Message); err != nil {
			h.log.Error("register: write invalid body error response failed", slog.Int("status", apiErr.Status), slog.String("code", string(apiErr.Code)), slog.String("error", err.Error()))
		}
		return
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		if err == nil {
			h.log.Warn("register invalid request body: extra json payload")
		} else {
			h.log.Warn("register invalid request body: extra json payload", slog.String("error", err.Error()))
		}
		apiErr := utils.MapError(utils.ErrInvalidRequest)
		if err := utils.WriteError(w, apiErr.Status, apiErr.Code, apiErr.Message); err != nil {
			h.log.Error("register: write extra payload error response failed", slog.Int("status", apiErr.Status), slog.String("code", string(apiErr.Code)), slog.String("error", err.Error()))
		}
		return
	}

	if err := utils.Validate(req); err != nil {
		h.log.Warn("register request validation failed", slog.String("login", req.Login), slog.String("error", err.Error()))
		apiErr := utils.MapError(utils.ErrInvalidRequest)
		if err := utils.WriteError(w, apiErr.Status, apiErr.Code, apiErr.Message); err != nil {
			h.log.Error("register: write validation error response failed", slog.Int("status", apiErr.Status), slog.String("code", string(apiErr.Code)), slog.String("error", err.Error()))
		}
		return
	}

	user, err := h.authService.Register(r.Context(), req.Login, req.Password, models.UserRoleUser)
	if err != nil {
		h.log.Warn("register failed", slog.String("login", req.Login), slog.String("error", err.Error()))
		apiErr := utils.MapError(err)
		if err := utils.WriteError(w, apiErr.Status, apiErr.Code, apiErr.Message); err != nil {
			h.log.Error("register: write service error response failed", slog.Int("status", apiErr.Status), slog.String("code", string(apiErr.Code)), slog.String("error", err.Error()))
		}
		return
	}

	var createdAt *time.Time
	if !user.CreatedAt.IsZero() {
		t := user.CreatedAt.UTC()
		createdAt = &t
	}

	h.log.Info("register completed", slog.String("user_id", user.ID.String()), slog.String("login", user.Login), slog.String("role", string(user.Role)))
	if err := utils.WriteJSON(w, http.StatusCreated, registerResponse{
		User: registerUserResponse{
			ID:        user.ID,
			Login:     user.Login,
			Role:      user.Role,
			CreatedAt: createdAt,
		},
	}); err != nil {
		h.log.Error("register: write success response failed", slog.String("user_id", user.ID.String()), slog.String("error", err.Error()))
	}
}
