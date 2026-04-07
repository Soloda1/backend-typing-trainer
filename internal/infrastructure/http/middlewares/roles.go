package middlewares

import (
	"log/slog"
	"net/http"

	"backend-typing-trainer/internal/domain/models"
	ports "backend-typing-trainer/internal/domain/ports/output/logger"
	"backend-typing-trainer/internal/utils"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func RequireRoles(log ports.Logger, roles ...models.UserRole) func(next http.Handler) http.Handler {
	allowed := make(map[models.UserRole]struct{}, len(roles))
	for _, role := range roles {
		allowed[role] = struct{}{}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			actor, ok := utils.ActorFromContext(r.Context())
			if !ok {
				apiErr := utils.MapError(utils.ErrUnauthorized)
				if err := utils.WriteError(w, apiErr.Status, apiErr.Code, apiErr.Message); err != nil {
					log.Error("roles middleware: write missing actor response failed", slog.String("error", err.Error()))
				}
				return
			}

			if _, exists := allowed[actor.Role]; !exists {
				apiErr := utils.MapError(utils.ErrForbidden)
				if err := utils.WriteError(w, apiErr.Status, apiErr.Code, apiErr.Message); err != nil {
					log.Error("roles middleware: write forbidden response failed", slog.String("role", string(actor.Role)), slog.String("error", err.Error()))
				}
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func RequireSelfOrRoles(log ports.Logger, userIDParam string, roles ...models.UserRole) func(next http.Handler) http.Handler {
	allowed := make(map[models.UserRole]struct{}, len(roles))
	for _, role := range roles {
		allowed[role] = struct{}{}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			actor, ok := utils.ActorFromContext(r.Context())
			if !ok {
				apiErr := utils.MapError(utils.ErrUnauthorized)
				if err := utils.WriteError(w, apiErr.Status, apiErr.Code, apiErr.Message); err != nil {
					log.Error("roles middleware: write missing actor response failed", slog.String("error", err.Error()))
				}
				return
			}

			targetID, err := uuid.Parse(chi.URLParam(r, userIDParam))
			if err != nil {
				apiErr := utils.MapError(utils.ErrInvalidRequest)
				if writeErr := utils.WriteError(w, apiErr.Status, apiErr.Code, apiErr.Message); writeErr != nil {
					log.Error("roles middleware: write invalid user id response failed", slog.String("param", userIDParam), slog.String("error", writeErr.Error()))
				}
				return
			}

			if actor.UserID == targetID {
				next.ServeHTTP(w, r)
				return
			}

			if _, exists := allowed[actor.Role]; !exists {
				apiErr := utils.MapError(utils.ErrForbidden)
				if err := utils.WriteError(w, apiErr.Status, apiErr.Code, apiErr.Message); err != nil {
					log.Error("roles middleware: write forbidden response failed", slog.String("role", string(actor.Role)), slog.String("error", err.Error()))
				}
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
