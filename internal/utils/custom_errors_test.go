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
			name:    "cannot cancel foreign booking",
			err:     ErrCannotCancelAnotherUsersBooking,
			status:  http.StatusForbidden,
			code:    ErrorCodeForbidden,
			message: "cannot cancel another user's booking",
		},
		{
			name:    "room not found",
			err:     ErrRoomNotFound,
			status:  http.StatusNotFound,
			code:    ErrorCodeRoomNotFound,
			message: "room not found",
		},
		{
			name:    "schedule not found",
			err:     ErrScheduleNotFound,
			status:  http.StatusNotFound,
			code:    ErrorCodeNotFound,
			message: "schedule not found",
		},
		{
			name:    "slot not found",
			err:     ErrSlotNotFound,
			status:  http.StatusNotFound,
			code:    ErrorCodeSlotNotFound,
			message: "slot not found",
		},
		{
			name:    "booking not found",
			err:     ErrBookingNotFound,
			status:  http.StatusNotFound,
			code:    ErrorCodeBookingNotFound,
			message: "booking not found",
		},
		{
			name:    "slot already booked",
			err:     ErrSlotAlreadyBooked,
			status:  http.StatusConflict,
			code:    ErrorCodeSlotAlreadyBooked,
			message: "slot is already booked",
		},
		{
			name:    "schedule exists",
			err:     ErrScheduleExists,
			status:  http.StatusConflict,
			code:    ErrorCodeScheduleExists,
			message: "schedule for this room already exists and cannot be changed",
		},
		{
			name:    "conference service unavailable",
			err:     ErrConferenceServiceUnavailable,
			status:  http.StatusInternalServerError,
			code:    ErrorCodeInternalError,
			message: "internal server error",
		},
		{
			name:    "wrapped conference service unavailable",
			err:     fmt.Errorf("wrapped: %w", ErrConferenceServiceUnavailable),
			status:  http.StatusInternalServerError,
			code:    ErrorCodeInternalError,
			message: "internal server error",
		},
		{
			name:    "conference service timeout",
			err:     ErrConferenceServiceTimeout,
			status:  http.StatusInternalServerError,
			code:    ErrorCodeInternalError,
			message: "internal server error",
		},
		{
			name:    "wrapped conference service timeout",
			err:     fmt.Errorf("wrapped: %w", ErrConferenceServiceTimeout),
			status:  http.StatusInternalServerError,
			code:    ErrorCodeInternalError,
			message: "internal server error",
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
