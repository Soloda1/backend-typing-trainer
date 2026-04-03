package auth

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	"backend-typing-trainer/internal/utils"
)

type loginRequest struct {
	Login    string `json:"login" validate:"required" example:"player_one"`
	Password string `json:"password" validate:"required" example:"secret"`
}

type LoginResponse struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"`
}

// Login godoc
// @Summary Авторизация по логину и паролю
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body loginRequest true "Login payload"
// @Success 200 {object} LoginResponse "Успешная авторизация"
// @Failure 401 {object} utils.ErrorResponse "Неверные учётные данные"
// @Failure 500 {object} utils.InternalErrorResponse "Внутренняя ошибка сервера"
// @Router /login [post]
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	h.log.Debug("login started")

	if h.authService == nil {
		h.log.Error("login failed: auth service is nil")
		if err := utils.WriteError(w, http.StatusInternalServerError, utils.ErrorCodeInternalError, "internal server error"); err != nil {
			h.log.Error("login: write 500 response failed", slog.String("error", err.Error()))
		}
		return
	}

	var req loginRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		h.log.Warn("login invalid request body", slog.String("error", err.Error()))
		apiErr := utils.MapError(utils.ErrUnauthorized)
		if err := utils.WriteError(w, apiErr.Status, apiErr.Code, apiErr.Message); err != nil {
			h.log.Error("login: write invalid body error response failed", slog.Int("status", apiErr.Status), slog.String("code", string(apiErr.Code)), slog.String("error", err.Error()))
		}
		return
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		if err == nil {
			h.log.Warn("login invalid request body: extra json payload")
		} else {
			h.log.Warn("login invalid request body: extra json payload", slog.String("error", err.Error()))
		}
		apiErr := utils.MapError(utils.ErrUnauthorized)
		if err := utils.WriteError(w, apiErr.Status, apiErr.Code, apiErr.Message); err != nil {
			h.log.Error("login: write extra payload error response failed", slog.Int("status", apiErr.Status), slog.String("code", string(apiErr.Code)), slog.String("error", err.Error()))
		}
		return
	}

	if err := utils.Validate(req); err != nil {
		h.log.Warn("login request validation failed", slog.String("login", req.Login), slog.String("error", err.Error()))
		apiErr := utils.MapError(utils.ErrUnauthorized)
		if err := utils.WriteError(w, apiErr.Status, apiErr.Code, apiErr.Message); err != nil {
			h.log.Error("login: write validation error response failed", slog.Int("status", apiErr.Status), slog.String("code", string(apiErr.Code)), slog.String("error", err.Error()))
		}
		return
	}

	token, err := h.authService.Login(r.Context(), req.Login, req.Password)
	if err != nil {
		h.log.Warn("login failed", slog.String("login", req.Login), slog.String("error", err.Error()))
		apiErr := utils.MapError(err)
		if err := utils.WriteError(w, apiErr.Status, apiErr.Code, apiErr.Message); err != nil {
			h.log.Error("login: write service error response failed", slog.Int("status", apiErr.Status), slog.String("code", string(apiErr.Code)), slog.String("error", err.Error()))
		}
		return
	}

	h.log.Info("login completed", slog.String("login", req.Login))
	if err := utils.WriteJSON(w, http.StatusOK, LoginResponse{Token: token}); err != nil {
		h.log.Error("login: write success response failed", slog.String("login", req.Login), slog.String("error", err.Error()))
	}
}
