package jwt

import (
	"testing"
	"time"

	jwtlib "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"backend-typing-trainer/internal/domain/models"
	infraLogger "backend-typing-trainer/internal/infrastructure/logger"
	"backend-typing-trainer/internal/utils"
)

const (
	testJWTSecret = "test-secret"
	testJWTIssuer = "booking-service"
)

func TestManager_NewTokenAndParseToken_Success(t *testing.T) {
	m := NewManager(testJWTSecret, time.Hour, testJWTIssuer, infraLogger.New("dev"))
	userID := uuid.MustParse("22222222-2222-2222-2222-222222222222")

	token, err := m.NewToken(userID, models.UserRoleUser)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	parsed, err := m.ParseToken(token)
	require.NoError(t, err)
	require.Equal(t, userID, parsed.UserID)
	require.Equal(t, models.UserRoleUser, parsed.Role)
}

func TestManager_ParseToken(t *testing.T) {
	m := NewManager(testJWTSecret, time.Hour, testJWTIssuer, infraLogger.New("dev"))
	userID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	secret := []byte(testJWTSecret)

	tests := []struct {
		name  string
		token string
	}{
		{
			name: "expired token",
			token: signToken(t, secret, claims{
				UserID: userID.String(),
				Role:   models.UserRoleAdmin,
				RegisteredClaims: jwtlib.RegisteredClaims{
					Issuer:    testJWTIssuer,
					IssuedAt:  jwtlib.NewNumericDate(time.Now().UTC().Add(-2 * time.Hour)),
					ExpiresAt: jwtlib.NewNumericDate(time.Now().UTC().Add(-time.Hour)),
				},
			}, jwtlib.SigningMethodHS256),
		},
		{
			name: "invalid issuer",
			token: signToken(t, secret, claims{
				UserID: userID.String(),
				Role:   models.UserRoleAdmin,
				RegisteredClaims: jwtlib.RegisteredClaims{
					Issuer:    "another-service",
					IssuedAt:  jwtlib.NewNumericDate(time.Now().UTC()),
					ExpiresAt: jwtlib.NewNumericDate(time.Now().UTC().Add(time.Hour)),
				},
			}, jwtlib.SigningMethodHS256),
		},
		{
			name: "invalid signing method",
			token: signToken(t, secret, claims{
				UserID: userID.String(),
				Role:   models.UserRoleAdmin,
				RegisteredClaims: jwtlib.RegisteredClaims{
					Issuer:    testJWTIssuer,
					IssuedAt:  jwtlib.NewNumericDate(time.Now().UTC()),
					ExpiresAt: jwtlib.NewNumericDate(time.Now().UTC().Add(time.Hour)),
				},
			}, jwtlib.SigningMethodHS512),
		},
		{
			name: "invalid role",
			token: signToken(t, secret, claims{
				UserID: userID.String(),
				Role:   models.UserRole("manager"),
				RegisteredClaims: jwtlib.RegisteredClaims{
					Issuer:    testJWTIssuer,
					IssuedAt:  jwtlib.NewNumericDate(time.Now().UTC()),
					ExpiresAt: jwtlib.NewNumericDate(time.Now().UTC().Add(time.Hour)),
				},
			}, jwtlib.SigningMethodHS256),
		},
		{
			name: "invalid user id",
			token: signToken(t, secret, claims{
				UserID: "not-uuid",
				Role:   models.UserRoleAdmin,
				RegisteredClaims: jwtlib.RegisteredClaims{
					Issuer:    testJWTIssuer,
					IssuedAt:  jwtlib.NewNumericDate(time.Now().UTC()),
					ExpiresAt: jwtlib.NewNumericDate(time.Now().UTC().Add(time.Hour)),
				},
			}, jwtlib.SigningMethodHS256),
		},
		{
			name:  "malformed token",
			token: "not-a-jwt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parsed, err := m.ParseToken(tt.token)
			require.ErrorIs(t, err, utils.ErrUnauthorized)
			require.Nil(t, parsed)
		})
	}
}

func TestIsAllowedRole(t *testing.T) {
	require.True(t, isAllowedRole(models.UserRoleAdmin))
	require.True(t, isAllowedRole(models.UserRoleUser))
	require.False(t, isAllowedRole(models.UserRole("manager")))
}

func signToken(t *testing.T, secret []byte, c claims, method jwtlib.SigningMethod) string {
	t.Helper()

	tok := jwtlib.NewWithClaims(method, c)
	signed, err := tok.SignedString(secret)
	require.NoError(t, err)

	return signed
}
