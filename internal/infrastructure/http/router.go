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
	exerciseshandler "backend-typing-trainer/internal/infrastructure/http/handlers/exercises"
	keyboardzoneshandler "backend-typing-trainer/internal/infrastructure/http/handlers/keyboard_zones"
	statisticshandler "backend-typing-trainer/internal/infrastructure/http/handlers/statistics"

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
	exercisesService        input.Exercises
	keyboardZonesService    input.KeyboardZones
	statisticsService       input.Statistics

	tokenManager jwtport.TokenManager
}

func NewRouter(log ports.Logger, authService input.Auth, difficultyLevelsService input.DifficultyLevels, exercisesService input.Exercises, keyboardZonesService input.KeyboardZones, statisticsService input.Statistics, tokenManager jwtport.TokenManager) *Router {
	return &Router{
		router:                  chi.NewRouter(),
		log:                     log,
		authService:             authService,
		difficultyLevelsService: difficultyLevelsService,
		exercisesService:        exercisesService,
		keyboardZonesService:    keyboardZonesService,
		statisticsService:       statisticsService,

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

		exercisesHandler := exerciseshandler.NewHandler(r.log, r.exercisesService)
		protected.Get("/exercises", exercisesHandler.List)
		protected.Get("/exercises/{id}", exercisesHandler.GetByID)

		keyboardZonesHandler := keyboardzoneshandler.NewHandler(r.log, r.keyboardZonesService)
		protected.Get("/keyboard-zones", keyboardZonesHandler.List)
		protected.Get("/keyboard-zones/{id}", keyboardZonesHandler.GetByID)
		protected.Get("/keyboard-zones/by-name/{name}", keyboardZonesHandler.GetByName)

		statisticsHandler := statisticshandler.NewHandler(r.log, r.statisticsService)
		protected.Post("/statistics", statisticsHandler.Create)
		protected.With(middlewares.RequireSelfOrRoles(r.log, "user_id", models.UserRoleAdmin)).Get("/statistics/users/{user_id}", statisticsHandler.ListByUserID)

		protected.Group(func(admin chi.Router) {
			admin.Use(middlewares.RequireRoles(r.log, models.UserRoleAdmin))
			admin.Post("/difficulty-levels", h.Create)
			admin.Patch("/difficulty-levels/{id}", h.Update)
			admin.Delete("/difficulty-levels/{id}", h.Delete)
			admin.Post("/exercises", exercisesHandler.Create)
			admin.Patch("/exercises/{id}", exercisesHandler.Update)
			admin.Delete("/exercises/{id}", exercisesHandler.Delete)
			admin.Get("/statistics/levels/{level_id}", statisticsHandler.ListByLevelID)
			admin.Get("/statistics/exercises/{exercise_id}", statisticsHandler.ListByExerciseID)
		})
	})
}

func (r *Router) GetRouter() *chi.Mux { return r.router }
