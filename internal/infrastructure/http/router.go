package http

import (
	"net/http"

	input "backend-typing-trainer/internal/domain/ports/input"
	jwtport "backend-typing-trainer/internal/domain/ports/output/jwt"
	ports "backend-typing-trainer/internal/domain/ports/output/logger"
	"backend-typing-trainer/internal/infrastructure/config"
	authhandler "backend-typing-trainer/internal/infrastructure/http/handlers/auth"

	"backend-typing-trainer/internal/infrastructure/http/middlewares"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	//httpSwagger "github.com/swaggo/http-swagger"
)

type Router struct {
	router      *chi.Mux
	log         ports.Logger
	authService input.Auth

	tokenManager jwtport.TokenManager
}

func NewRouter(log ports.Logger, authService input.Auth, tokenManager jwtport.TokenManager) *Router {
	return &Router{
		router:      chi.NewRouter(),
		log:         log,
		authService: authService,

		tokenManager: tokenManager,
	}
}

func (r *Router) Setup(cfg *config.Config) {
	r.router.Use(chiMiddleware.RequestID)
	r.router.Use(chiMiddleware.RealIP)
	r.router.Use(middlewares.RequestLoggerMiddleware(r.log))
	r.router.Use(chiMiddleware.Recoverer)
	r.router.Use(chiMiddleware.Timeout(cfg.HTTPServer.RequestTimeout))

	//r.router.Get("/swagger/*", httpSwagger.WrapHandler)
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
	})
}

func (r *Router) GetRouter() *chi.Mux { return r.router }
