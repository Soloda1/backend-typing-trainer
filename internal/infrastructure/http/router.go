package http

import (
	"net/http"

	"backend-typing-trainer/internal/domain/models"
	input "backend-typing-trainer/internal/domain/ports/input"
	jwtport "backend-typing-trainer/internal/domain/ports/output/jwt"
	ports "backend-typing-trainer/internal/domain/ports/output/logger"
	"backend-typing-trainer/internal/infrastructure/config"
	authhandler "backend-typing-trainer/internal/infrastructure/http/handlers/auth"
	difficultylevelshandler "backend-typing-trainer/internal/infrastructure/http/handlers/difficulty_levels"
	keyboardzoneshandler "backend-typing-trainer/internal/infrastructure/http/handlers/keyboard_zones"

	"backend-typing-trainer/internal/infrastructure/http/middlewares"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Router struct {
	router                  *chi.Mux
	log                     ports.Logger
	authService             input.Auth
	difficultyLevelsService input.DifficultyLevels
	keyboardZonesService    input.KeyboardZones

	tokenManager jwtport.TokenManager
}

func NewRouter(log ports.Logger, authService input.Auth, difficultyLevelsService input.DifficultyLevels, keyboardZonesService input.KeyboardZones, tokenManager jwtport.TokenManager) *Router {
	return &Router{
		router:                  chi.NewRouter(),
		log:                     log,
		authService:             authService,
		difficultyLevelsService: difficultyLevelsService,
		keyboardZonesService:    keyboardZonesService,

		tokenManager: tokenManager,
	}
}

func (r *Router) Setup(cfg *config.Config) {
	r.router.Use(chiMiddleware.RequestID)
	r.router.Use(chiMiddleware.RealIP)
	r.router.Use(middlewares.RequestLoggerMiddleware(r.log))
	r.router.Use(chiMiddleware.Recoverer)
	r.router.Use(chiMiddleware.Timeout(cfg.HTTPServer.RequestTimeout))

	r.router.Get("/swagger/*", httpSwagger.WrapHandler)
	r.router.Mount("/", r.setupAuthRoutes())
	r.setupProtectedRoutes()
}

func (r *Router) setupAuthRoutes() http.Handler {
	h := authhandler.NewHandler(r.log, r.authService)
	sub := chi.NewRouter()
	sub.Post("/register", h.Register)
	sub.Post("/login", h.Login)
	return sub
}

func (r *Router) setupProtectedRoutes() {
	r.router.Group(func(protected chi.Router) {
		protected.Use(middlewares.AuthMiddleware(r.tokenManager, r.log))

		h := difficultylevelshandler.NewHandler(r.log, r.difficultyLevelsService)
		protected.Get("/difficulty-levels", h.List)
		protected.Get("/difficulty-levels/{id}", h.GetByID)

		keyboardZonesHandler := keyboardzoneshandler.NewHandler(r.log, r.keyboardZonesService)
		protected.Get("/keyboard-zones", keyboardZonesHandler.List)
		protected.Get("/keyboard-zones/{id}", keyboardZonesHandler.GetByID)
		protected.Get("/keyboard-zones/by-name/{name}", keyboardZonesHandler.GetByName)

		protected.Group(func(admin chi.Router) {
			admin.Use(middlewares.RequireRoles(r.log, models.UserRoleAdmin))
			admin.Post("/difficulty-levels", h.Create)
			admin.Patch("/difficulty-levels/{id}", h.Update)
			admin.Delete("/difficulty-levels/{id}", h.Delete)
		})
	})
}

func (r *Router) GetRouter() *chi.Mux { return r.router }
