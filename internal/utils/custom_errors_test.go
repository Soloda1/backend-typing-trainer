package utils

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMapError(t *testing.T) {
	tests := []struct {
		name    string
		err     error
		status  int
		code    ErrorCode
		message string
	}{
		{
			name:    "nil error",
			err:     nil,
			status:  http.StatusInternalServerError,
			code:    ErrorCodeInternalError,
			message: "internal server error",
		},
		{
			name:    "invalid request",
			err:     ErrInvalidRequest,
			status:  http.StatusBadRequest,
			code:    ErrorCodeInvalidRequest,
			message: "invalid request",
		},
		{
			name:    "wrapped invalid request",
			err:     fmt.Errorf("wrapped: %w", ErrInvalidRequest),
			status:  http.StatusBadRequest,
			code:    ErrorCodeInvalidRequest,
			message: "invalid request",
		},
		{
			name:    "unauthorized",
			err:     ErrUnauthorized,
			status:  http.StatusUnauthorized,
			code:    ErrorCodeUnauthorized,
			message: "unauthorized",
		},
		{
			name:    "forbidden",
			err:     ErrForbidden,
			status:  http.StatusForbidden,
			code:    ErrorCodeForbidden,
			message: "forbidden",
		},
		{
			name:    "user login exists",
			err:     ErrUserLoginExists,
			status:  http.StatusBadRequest,
			code:    ErrorCodeInvalidRequest,
			message: "invalid request",
		},
		{
			name:    "unknown error",
			err:     errors.New("boom"),
			status:  http.StatusInternalServerError,
			code:    ErrorCodeInternalError,
			message: "internal server error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MapError(tt.err)
			require.Equal(t, tt.status, got.Status)
			require.Equal(t, tt.code, got.Code)
			require.Equal(t, tt.message, got.Message)
		})
	}
}
