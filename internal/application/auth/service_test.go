package auth_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	authapp "backend-typing-trainer/internal/application/auth"
	"backend-typing-trainer/internal/domain/models"
	"backend-typing-trainer/internal/utils"
	"backend-typing-trainer/mocks"
)

func newServiceWithMocks(t *testing.T, tm *mocks.TokenManager, repo *mocks.UserRepository) *authapp.Service {
	t.Helper()

	log := mocks.NewLogger(t)
	log.EXPECT().With(mock.Anything).Return(log).Once()
	log.EXPECT().Warn(mock.Anything, mock.Anything).Maybe()
	log.EXPECT().Warn(mock.Anything, mock.Anything, mock.Anything).Maybe()
	log.EXPECT().Error(mock.Anything, mock.Anything).Maybe()
	log.EXPECT().Error(mock.Anything, mock.Anything, mock.Anything).Maybe()

	return authapp.NewService(tm, repo, log)
}

func TestService_Register(t *testing.T) {
	tests := []struct {
		name              string
		login             string
		password          string
		role              models.UserRole
		createErr         error
		wantErr           error
		wantCreateCalls   int
		wantLogin         string
		wantRole          models.UserRole
		checkPasswordHash bool
	}{
		{
			name:            "empty login",
			login:           "   ",
			password:        "secret",
			role:            models.UserRoleUser,
			wantErr:         utils.ErrInvalidRequest,
			wantCreateCalls: 0,
		},
		{
			name:            "empty password",
			login:           "player_one",
			password:        "   ",
			role:            models.UserRoleUser,
			wantErr:         utils.ErrInvalidRequest,
			wantCreateCalls: 0,
		},
		{
			name:            "invalid role",
			login:           "player_one",
			password:        "secret",
			role:            models.UserRole("guest"),
			wantErr:         utils.ErrInvalidRequest,
			wantCreateCalls: 0,
		},
		{
			name:            "duplicate login",
			login:           "player_one",
			password:        "secret",
			role:            models.UserRoleUser,
			createErr:       utils.ErrUserLoginExists,
			wantErr:         utils.ErrUserLoginExists,
			wantCreateCalls: 1,
			wantLogin:       "player_one",
			wantRole:        models.UserRoleUser,
		},
		{
			name:              "success trims login and hashes password",
			login:             "  player_one  ",
			password:          "secret",
			role:              models.UserRoleAdmin,
			wantCreateCalls:   1,
			wantLogin:         "player_one",
			wantRole:          models.UserRoleAdmin,
			checkPasswordHash: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mocks.NewUserRepository(t)
			tm := mocks.NewTokenManager(t)
			svc := newServiceWithMocks(t, tm, repo)

			var createdUser *models.User
			if tt.wantCreateCalls > 0 {
				repo.EXPECT().Create(mock.Anything, mock.Anything).
					Run(func(_ context.Context, user *models.User) {
						createdUser = user
						if tt.createErr == nil {
							user.ID = uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa")
						}
					}).
					Return(tt.createErr).
					Once()
			}

			user, err := svc.Register(context.Background(), tt.login, tt.password, tt.role)

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("expected error %v, got %v", tt.wantErr, err)
				}
			} else if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			repo.AssertNumberOfCalls(t, "Create", tt.wantCreateCalls)

			if tt.wantCreateCalls > 0 {
				if createdUser == nil {
					t.Fatal("expected created user to be captured")
				}
				if createdUser.Login != tt.wantLogin {
					t.Fatalf("expected login %q, got %q", tt.wantLogin, createdUser.Login)
				}
				if createdUser.Role != tt.wantRole {
					t.Fatalf("expected role %q, got %q", tt.wantRole, createdUser.Role)
				}
				if createdUser.PasswordHash == nil || *createdUser.PasswordHash == "" {
					t.Fatal("expected password hash to be set")
				}
				if tt.checkPasswordHash && !utils.CheckPasswordHash(tt.password, *createdUser.PasswordHash) {
					t.Fatal("stored hash does not match original password")
				}
			}

			if tt.wantErr == nil {
				if user == nil {
					t.Fatal("expected returned user")
				}
				if user.Login != tt.wantLogin {
					t.Fatalf("expected returned login %q, got %q", tt.wantLogin, user.Login)
				}
			}
		})
	}
}

func TestService_Login(t *testing.T) {
	passwordHash, err := utils.HashPassword("secret")
	if err != nil {
		t.Fatalf("hash password in test setup: %v", err)
	}
	errDBTimeout := errors.New("db timeout")
	errSignFailed := errors.New("sign failed")

	tests := []struct {
		name            string
		login           string
		password        string
		user            *models.User
		getByLoginErr   error
		tokenResult     string
		tokenErr        error
		wantErr         error
		wantToken       string
		wantRepoCalls   int
		wantTokenCalls  int
		wantLookupLogin string
	}{
		{
			name:          "empty login",
			login:         "   ",
			password:      "secret",
			wantErr:       utils.ErrUnauthorized,
			wantRepoCalls: 0,
		},
		{
			name:          "empty password",
			login:         "player_one",
			password:      "   ",
			wantErr:       utils.ErrUnauthorized,
			wantRepoCalls: 0,
		},
		{
			name:            "user not found",
			login:           "player_one",
			password:        "secret",
			getByLoginErr:   utils.ErrUserNotFound,
			wantErr:         utils.ErrUnauthorized,
			wantRepoCalls:   1,
			wantLookupLogin: "player_one",
		},
		{
			name:            "repository infra error",
			login:           "player_one",
			password:        "secret",
			getByLoginErr:   errDBTimeout,
			wantErr:         errDBTimeout,
			wantRepoCalls:   1,
			wantLookupLogin: "player_one",
		},
		{
			name:     "missing password hash",
			login:    "player_one",
			password: "secret",
			user: &models.User{
				ID:    uuid.MustParse("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"),
				Login: "player_one",
				Role:  models.UserRoleUser,
			},
			wantErr:         utils.ErrUnauthorized,
			wantRepoCalls:   1,
			wantLookupLogin: "player_one",
		},
		{
			name:     "wrong password",
			login:    "player_one",
			password: "wrong",
			user: &models.User{
				ID:           uuid.MustParse("cccccccc-cccc-cccc-cccc-cccccccccccc"),
				Login:        "player_one",
				Role:         models.UserRoleUser,
				PasswordHash: &passwordHash,
			},
			wantErr:         utils.ErrUnauthorized,
			wantRepoCalls:   1,
			wantLookupLogin: "player_one",
		},
		{
			name:     "token manager error",
			login:    "player_one",
			password: "secret",
			user: &models.User{
				ID:           uuid.MustParse("dddddddd-dddd-dddd-dddd-dddddddddddd"),
				Login:        "player_one",
				Role:         models.UserRoleAdmin,
				PasswordHash: &passwordHash,
			},
			tokenErr:        errSignFailed,
			wantErr:         errSignFailed,
			wantRepoCalls:   1,
			wantTokenCalls:  1,
			wantLookupLogin: "player_one",
		},
		{
			name:     "success trims login",
			login:    "  player_one  ",
			password: "secret",
			user: &models.User{
				ID:           uuid.MustParse("eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee"),
				Login:        "player_one",
				Role:         models.UserRoleAdmin,
				PasswordHash: &passwordHash,
			},
			tokenResult:     "jwt-token",
			wantToken:       "jwt-token",
			wantRepoCalls:   1,
			wantTokenCalls:  1,
			wantLookupLogin: "player_one",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mocks.NewUserRepository(t)
			tm := mocks.NewTokenManager(t)
			svc := newServiceWithMocks(t, tm, repo)

			if tt.wantRepoCalls > 0 {
				repo.EXPECT().GetByLogin(mock.Anything, tt.wantLookupLogin).Return(tt.user, tt.getByLoginErr).Once()
			}

			if tt.wantTokenCalls > 0 {
				tm.EXPECT().NewToken(tt.user.ID, tt.user.Role).Return(tt.tokenResult, tt.tokenErr).Once()
			}

			gotToken, err := svc.Login(context.Background(), tt.login, tt.password)

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("expected error %v, got %v", tt.wantErr, err)
				}
			} else if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if gotToken != tt.wantToken {
				t.Fatalf("expected token %q, got %q", tt.wantToken, gotToken)
			}

			repo.AssertNumberOfCalls(t, "GetByLogin", tt.wantRepoCalls)
			tm.AssertNumberOfCalls(t, "NewToken", tt.wantTokenCalls)
		})
	}
}
