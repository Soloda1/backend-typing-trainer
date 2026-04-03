package utils

import (
	"errors"
	"net/http"
)

type ErrorCode string

const (
	ErrorCodeInvalidRequest ErrorCode = "INVALID_REQUEST"
	ErrorCodeUnauthorized   ErrorCode = "UNAUTHORIZED"
	ErrorCodeNotFound       ErrorCode = "NOT_FOUND"
	ErrorCodeForbidden      ErrorCode = "FORBIDDEN"
	ErrorCodeInternalError  ErrorCode = "INTERNAL_ERROR"
)

var (
	ErrInvalidRequest  = errors.New("invalid request")
	ErrUnauthorized    = errors.New("unauthorized")
	ErrForbidden       = errors.New("forbidden")
	ErrNotFound        = errors.New("not found")
	ErrUserNotFound    = errors.New("user not found")
	ErrUserLoginExists = errors.New("user login already exists")
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
	case errors.Is(err, ErrUserNotFound):
		return APIError{Status: http.StatusNotFound, Code: ErrorCodeNotFound, Message: "user not found"}
	case errors.Is(err, ErrNotFound):
		return APIError{Status: http.StatusNotFound, Code: ErrorCodeNotFound, Message: "not found"}
	case errors.Is(err, ErrUserLoginExists):
		return APIError{Status: http.StatusBadRequest, Code: ErrorCodeInvalidRequest, Message: "invalid request"}
	default:
		return APIError{Status: http.StatusInternalServerError, Code: ErrorCodeInternalError, Message: "internal server error"}
	}
}
