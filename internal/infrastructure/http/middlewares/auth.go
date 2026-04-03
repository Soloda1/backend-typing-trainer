package middlewares

import (
	"log/slog"
	"net/http"
	"strings"

	input "backend-typing-trainer/internal/domain/ports/input"
	jwtport "backend-typing-trainer/internal/domain/ports/output/jwt"
	ports "backend-typing-trainer/internal/domain/ports/output/logger"
	"backend-typing-trainer/internal/utils"
)

func AuthMiddleware(tokenManager jwtport.TokenManager, log ports.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if tokenManager == nil {
				log.Error("auth middleware misconfigured: token manager is nil")
				if err := utils.WriteError(w, http.StatusInternalServerError, utils.ErrorCodeInternalError, "internal server error"); err != nil {
					log.Error("auth middleware: write 500 response failed", slog.String("error", err.Error()))
				}
				return
			}

			authHeader := strings.TrimSpace(r.Header.Get("Authorization"))
			token, ok := extractBearerToken(authHeader)
			if !ok {
				apiErr := utils.MapError(utils.ErrUnauthorized)
				if err := utils.WriteError(w, apiErr.Status, apiErr.Code, apiErr.Message); err != nil {
					log.Error("auth middleware: write unauthorized response failed", slog.String("error", err.Error()))
				}
				return
			}

			claims, err := tokenManager.ParseToken(token)
			if err != nil {
				apiErr := utils.MapError(err)
				if err := utils.WriteError(w, apiErr.Status, apiErr.Code, apiErr.Message); err != nil {
					log.Error("auth middleware: write parse token error response failed", slog.String("error", err.Error()))
				}
				return
			}

			actor := input.Actor{
				UserID: claims.UserID,
				Role:   claims.Role,
			}

			next.ServeHTTP(w, r.WithContext(utils.WithActor(r.Context(), actor)))
		})
	}
}

func extractBearerToken(header string) (string, bool) {
	if header == "" {
		return "", false
	}

	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 {
		return "", false
	}

	if !strings.EqualFold(parts[0], "Bearer") {
		return "", false
	}

	token := strings.TrimSpace(parts[1])
	if token == "" {
		return "", false
	}

	return token, true
}
