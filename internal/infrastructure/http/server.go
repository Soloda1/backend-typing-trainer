package http

import (
	"context"

	input "backend-typing-trainer/internal/domain/ports/input"
	jwtport "backend-typing-trainer/internal/domain/ports/output/jwt"
	ports "backend-typing-trainer/internal/domain/ports/output/logger"
	"backend-typing-trainer/internal/infrastructure/config"
	"log/slog"
	"net/http"
)

type Server struct {
	address                 string
	log                     ports.Logger
	authService             input.Auth
	difficultyLevelsService input.DifficultyLevels
	exercisesService        input.Exercises
	keyboardZonesService    input.KeyboardZones
	statisticsService       input.Statistics
	tokenManager            jwtport.TokenManager
	router                  *Router
	server                  *http.Server
}

func NewServer(address string, log ports.Logger, authService input.Auth, difficultyLevelsService input.DifficultyLevels, exercisesService input.Exercises, keyboardZonesService input.KeyboardZones, statisticsService input.Statistics, tokenManager jwtport.TokenManager) *Server {
	return &Server{
		address:                 address,
		log:                     log,
		authService:             authService,
		difficultyLevelsService: difficultyLevelsService,
		exercisesService:        exercisesService,
		keyboardZonesService:    keyboardZonesService,
		statisticsService:       statisticsService,

		tokenManager: tokenManager,
	}
}

func (s *Server) Run(cfg *config.Config) error {
	s.router = NewRouter(s.log, s.authService, s.difficultyLevelsService, s.exercisesService, s.keyboardZonesService, s.statisticsService, s.tokenManager)
	s.router.Setup(cfg)

	s.server = &http.Server{
		Addr:         s.address,
		Handler:      s.router.GetRouter(),
		ReadTimeout:  cfg.HTTPServer.ReadTimeout,
		WriteTimeout: cfg.HTTPServer.WriteTimeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	s.log.Info("starting http server", slog.String("address", s.address))
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	if s.server == nil {
		return nil
	}
	return s.server.Shutdown(ctx)
}
