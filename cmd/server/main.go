// Package main Typing Trainer Service API.
//
// @title Typing Trainer Service API
// @version 1.0.0
// @description Backend API для сервиса клавиатурного тренажера.
// @BasePath /
// @schemes http
//
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description JWT token в формате `Bearer <token>`.
package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "backend-typing-trainer/docs"

	"github.com/jackc/pgx/v5/pgxpool"

	authapp "backend-typing-trainer/internal/application/auth"
	difficultylevelsapp "backend-typing-trainer/internal/application/difficulty_levels"
	keyboardzonesapp "backend-typing-trainer/internal/application/keyboard_zones"

	jwtmanager "backend-typing-trainer/internal/infrastructure/auth/jwt"
	"backend-typing-trainer/internal/infrastructure/config"
	httpserver "backend-typing-trainer/internal/infrastructure/http"
	"backend-typing-trainer/internal/infrastructure/logger"
	difficultylevelsrepo "backend-typing-trainer/internal/infrastructure/persistence/postgres/difficulty_levels"
	keyboardzonesrepo "backend-typing-trainer/internal/infrastructure/persistence/postgres/keyboard_zones"
	usersrepo "backend-typing-trainer/internal/infrastructure/persistence/postgres/users"
)

func main() {
	cfg := config.MustLoad()
	log := logger.New(cfg.Env)

	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.Database.Username,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.DbName,
	)

	dbCtx, dbCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer dbCancel()

	dbPool, err := pgxpool.New(dbCtx, dsn)
	if err != nil {
		log.Error("failed to connect database", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer dbPool.Close()

	if err := dbPool.Ping(dbCtx); err != nil {
		log.Error("failed to ping database", slog.String("error", err.Error()))
		os.Exit(1)
	}

	tokenManager := jwtmanager.NewManager(cfg.JWT.Secret, cfg.JWT.TTL, cfg.JWT.Issuer, log)
	usersRepository := usersrepo.NewRepository(dbPool, log)
	difficultyLevelsRepository := difficultylevelsrepo.NewRepository(dbPool, log)
	keyboardZonesRepository := keyboardzonesrepo.NewRepository(dbPool, log)

	authService := authapp.NewService(tokenManager, usersRepository, log)
	difficultyLevelsService := difficultylevelsapp.NewService(difficultyLevelsRepository, log)
	keyboardZonesService := keyboardzonesapp.NewService(keyboardZonesRepository, log)

	address := fmt.Sprintf("%s:%d", cfg.HTTPServer.Address, cfg.HTTPServer.Port)
	server := httpserver.NewServer(address, log, authService, difficultyLevelsService, keyboardZonesService, tokenManager)

	serverErr := make(chan error, 1)
	go func() {
		serverErr <- server.Run(cfg)
	}()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	select {
	case sig := <-signals:
		log.Info("shutdown signal received", slog.String("signal", sig.String()))
	case err := <-serverErr:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("http server failed", slog.String("error", err.Error()))
			os.Exit(1)
		}
		return
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Error("http server shutdown failed", slog.String("error", err.Error()))
		os.Exit(1)
	}

	if err := <-serverErr; err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Error("http server stopped with error", slog.String("error", err.Error()))
		os.Exit(1)
	}

	log.Info("server stopped")
}
