package utils

import (
	"errors"
	"net/http"
)

type ErrorCode string

const (
	ErrorCodeInvalidRequest    ErrorCode = "INVALID_REQUEST"
	ErrorCodeUnauthorized      ErrorCode = "UNAUTHORIZED"
	ErrorCodeNotFound          ErrorCode = "NOT_FOUND"
	ErrorCodeRoomNotFound      ErrorCode = "ROOM_NOT_FOUND"
	ErrorCodeSlotNotFound      ErrorCode = "SLOT_NOT_FOUND"
	ErrorCodeSlotAlreadyBooked ErrorCode = "SLOT_ALREADY_BOOKED"
	ErrorCodeBookingNotFound   ErrorCode = "BOOKING_NOT_FOUND"
	ErrorCodeForbidden         ErrorCode = "FORBIDDEN"
	ErrorCodeScheduleExists    ErrorCode = "SCHEDULE_EXISTS"
	ErrorCodeInternalError     ErrorCode = "INTERNAL_ERROR"
)

var (
	ErrInvalidRequest = errors.New("invalid request")
	ErrUnauthorized   = errors.New("unauthorized")
	ErrForbidden      = errors.New("forbidden")
	ErrNotFound       = errors.New("not found")

	ErrRoomNotFound     = errors.New("room not found")
	ErrScheduleNotFound = errors.New("schedule not found")
	ErrSlotNotFound     = errors.New("slot not found")
	ErrUserNotFound     = errors.New("user not found")

	ErrSlotAlreadyBooked = errors.New("slot already booked")
	ErrBookingNotFound   = errors.New("booking not found")
	ErrScheduleExists    = errors.New("schedule already exists")
	ErrUserLoginExists   = errors.New("user login already exists")

	ErrConferenceServiceUnavailable = errors.New("conference service unavailable")
	ErrConferenceServiceTimeout     = errors.New("conference service timeout")

	ErrCannotCancelAnotherUsersBooking = errors.New("cannot cancel another user's booking")
)

type APIError struct {
	Status  int
	Code    ErrorCode
	Message string
}

func MapError(err error) APIError {
	switch {
	case err == nil:
		return APIError{Status: http.StatusInternalServerError, Code: ErrorCodeInternalError, Message: "internal server error"}
	case errors.Is(err, ErrInvalidRequest):
		return APIError{Status: http.StatusBadRequest, Code: ErrorCodeInvalidRequest, Message: "invalid request"}
	case errors.Is(err, ErrUnauthorized):
		return APIError{Status: http.StatusUnauthorized, Code: ErrorCodeUnauthorized, Message: "unauthorized"}
	case errors.Is(err, ErrForbidden):
		return APIError{Status: http.StatusForbidden, Code: ErrorCodeForbidden, Message: "forbidden"}
	case errors.Is(err, ErrCannotCancelAnotherUsersBooking):
		return APIError{Status: http.StatusForbidden, Code: ErrorCodeForbidden, Message: "cannot cancel another user's booking"}
	case errors.Is(err, ErrRoomNotFound):
		return APIError{Status: http.StatusNotFound, Code: ErrorCodeRoomNotFound, Message: "room not found"}
	case errors.Is(err, ErrScheduleNotFound):
		return APIError{Status: http.StatusNotFound, Code: ErrorCodeNotFound, Message: "schedule not found"}
	case errors.Is(err, ErrUserNotFound):
		return APIError{Status: http.StatusNotFound, Code: ErrorCodeNotFound, Message: "user not found"}
	case errors.Is(err, ErrSlotNotFound):
		return APIError{Status: http.StatusNotFound, Code: ErrorCodeSlotNotFound, Message: "slot not found"}
	case errors.Is(err, ErrBookingNotFound):
		return APIError{Status: http.StatusNotFound, Code: ErrorCodeBookingNotFound, Message: "booking not found"}
	case errors.Is(err, ErrNotFound):
		return APIError{Status: http.StatusNotFound, Code: ErrorCodeNotFound, Message: "not found"}
	case errors.Is(err, ErrSlotAlreadyBooked):
		return APIError{Status: http.StatusConflict, Code: ErrorCodeSlotAlreadyBooked, Message: "slot is already booked"}
	case errors.Is(err, ErrScheduleExists):
		return APIError{Status: http.StatusConflict, Code: ErrorCodeScheduleExists, Message: "schedule for this room already exists and cannot be changed"}
	case errors.Is(err, ErrConferenceServiceUnavailable):
		return APIError{Status: http.StatusInternalServerError, Code: ErrorCodeInternalError, Message: "internal server error"}
	case errors.Is(err, ErrConferenceServiceTimeout):
		return APIError{Status: http.StatusInternalServerError, Code: ErrorCodeInternalError, Message: "internal server error"}
	case errors.Is(err, ErrUserLoginExists):
		return APIError{Status: http.StatusBadRequest, Code: ErrorCodeInvalidRequest, Message: "invalid request"}
	default:
		return APIError{Status: http.StatusInternalServerError, Code: ErrorCodeInternalError, Message: "internal server error"}
	}
}
