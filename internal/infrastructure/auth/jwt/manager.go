package jwt

import (
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"backend-typing-trainer/internal/domain/models"
	jwtport "backend-typing-trainer/internal/domain/ports/output/jwt"
	"backend-typing-trainer/internal/domain/ports/output/logger"
	"backend-typing-trainer/internal/utils"
)

type Manager struct {
	secret []byte
	ttl    time.Duration
	issuer string
	log    logger.Logger
}

type claims struct {
	UserID string          `json:"user_id"`
	Role   models.UserRole `json:"role"`
	jwt.RegisteredClaims
}

func NewManager(secret string, ttl time.Duration, issuer string, log logger.Logger) jwtport.TokenManager {
	return &Manager{
		secret: []byte(secret),
		ttl:    ttl,
		issuer: issuer,
		log:    log.With(slog.String("component", "jwt_manager")),
	}
}

func (m *Manager) NewToken(userID uuid.UUID, role models.UserRole) (string, error) {
	now := time.Now().UTC()

	tokenClaims := claims{
		UserID: userID.String(),
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    m.issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(m.ttl)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims)
	signedToken, err := token.SignedString(m.secret)
	if err != nil {
		m.log.Error("sign token failed", slog.String("error", err.Error()))
		return "", fmt.Errorf("sign token: %w", err)
	}

	return signedToken, nil
}

func (m *Manager) ParseToken(token string) (*jwtport.TokenClaims, error) {
	parsedToken, err := jwt.ParseWithClaims(token, &claims{}, func(parsedToken *jwt.Token) (any, error) {
		if parsedToken.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, fmt.Errorf("unexpected signing method: %v", parsedToken.Header["alg"])
		}
		return m.secret, nil
	},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithIssuer(m.issuer),
	)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			m.log.Debug("token expired")
			return nil, utils.ErrUnauthorized
		}
		m.log.Warn("parse token failed", slog.String("error", err.Error()))
		return nil, utils.ErrUnauthorized
	}

	parsedClaims, ok := parsedToken.Claims.(*claims)
	if !ok || !parsedToken.Valid {
		m.log.Warn("invalid token claims")
		return nil, utils.ErrUnauthorized
	}

	userID, err := uuid.Parse(parsedClaims.UserID)
	if err != nil {
		m.log.Warn("invalid user_id in token", slog.String("error", err.Error()))
		return nil, utils.ErrUnauthorized
	}

	if !isAllowedRole(parsedClaims.Role) {
		m.log.Warn("invalid role in token", slog.String("role", string(parsedClaims.Role)))
		return nil, utils.ErrUnauthorized
	}

	return &jwtport.TokenClaims{
		UserID: userID,
		Role:   parsedClaims.Role,
	}, nil
}

func isAllowedRole(role models.UserRole) bool {
	return role == models.UserRoleAdmin || role == models.UserRoleUser
}
