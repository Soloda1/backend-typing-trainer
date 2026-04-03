package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"backend-typing-trainer/internal/domain/models"
	jwtport "backend-typing-trainer/internal/domain/ports/output/jwt"
	ports "backend-typing-trainer/internal/domain/ports/output/logger"
	usersport "backend-typing-trainer/internal/domain/ports/output/users"
	"backend-typing-trainer/internal/utils"
)

type Service struct {
	tokenManager jwtport.TokenManager
	usersRepo    usersport.Repository
	log          ports.Logger
}

func NewService(tokenManager jwtport.TokenManager, usersRepo usersport.Repository, log ports.Logger) *Service {
	return &Service{
		tokenManager: tokenManager,
		usersRepo:    usersRepo,
		log:          log.With(slog.String("component", "application_auth")),
	}
}

func (s *Service) Login(ctx context.Context, login, password string) (string, error) {
	login = strings.TrimSpace(login)
	if !isValidLogin(login) || strings.TrimSpace(password) == "" {
		s.log.Warn("login rejected: invalid credentials format", slog.String("login", login))
		return "", utils.ErrUnauthorized
	}

	user, err := s.usersRepo.GetByLogin(ctx, login)
	if err != nil {
		if errors.Is(err, utils.ErrUserNotFound) {
			s.log.Warn("login rejected: invalid credentials", slog.String("login", login))
			return "", utils.ErrUnauthorized
		}
		s.log.Error("login failed: get user by login", slog.String("login", login), slog.String("error", err.Error()))
		return "", fmt.Errorf("login get user by login: %w", err)
	}

	if user.PasswordHash == nil || !utils.CheckPasswordHash(password, *user.PasswordHash) {
		s.log.Warn("login rejected: invalid credentials", slog.String("login", login))
		return "", utils.ErrUnauthorized
	}

	token, err := s.tokenManager.NewToken(user.ID, user.Role)
	if err != nil {
		s.log.Error("login failed: generate token", slog.String("user_id", user.ID.String()), slog.String("error", err.Error()))
		return "", fmt.Errorf("login generate token: %w", err)
	}

	return token, nil
}

func (s *Service) Register(ctx context.Context, login, password string, role models.UserRole) (*models.User, error) {
	login = strings.TrimSpace(login)
	if !isValidLogin(login) || strings.TrimSpace(password) == "" || !isAllowedRole(role) {
		s.log.Warn("register rejected: invalid input", slog.String("login", login), slog.String("role", string(role)))
		return nil, utils.ErrInvalidRequest
	}

	passwordHash, err := utils.HashPassword(password)
	if err != nil {
		s.log.Error("register failed: hash password", slog.String("error", err.Error()))
		return nil, fmt.Errorf("register hash password: %w", err)
	}

	user := &models.User{
		Login:        login,
		Role:         role,
		PasswordHash: &passwordHash,
	}

	if err := s.usersRepo.Create(ctx, user); err != nil {
		s.log.Warn("register failed", slog.String("login", login), slog.String("error", err.Error()))
		return nil, err
	}

	return user, nil
}

func isAllowedRole(role models.UserRole) bool {
	return role == models.UserRoleAdmin || role == models.UserRoleUser
}

func isValidLogin(login string) bool {
	return login != ""
}
