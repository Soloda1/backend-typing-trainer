package utils

import (
	"encoding/json"
	"net/http"
)

type ErrorDetails struct {
	Code    ErrorCode `json:"code" example:"INVALID_REQUEST"`
	Message string    `json:"message" example:"invalid request"`
}

type ErrorResponse struct {
	Error ErrorDetails `json:"error"`
}

type InternalErrorDetails struct {
	Code    ErrorCode `json:"code" example:"INTERNAL_ERROR"`
	Message string    `json:"message" example:"internal server error"`
}

type InternalErrorResponse struct {
	Error InternalErrorDetails `json:"error"`
}

func HTTPCodeConverter(status int) ErrorCode {
	switch status {
	case http.StatusBadRequest:
		return ErrorCodeInvalidRequest
	case http.StatusUnauthorized:
		return ErrorCodeUnauthorized
	case http.StatusForbidden:
		return ErrorCodeForbidden
	case http.StatusNotFound:
		return ErrorCodeNotFound
	default:
		return ErrorCodeInternalError
	}
}

func WriteJSON(w http.ResponseWriter, status int, payload any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(payload)
}

func WriteError(w http.ResponseWriter, status int, code ErrorCode, message string) error {
	if code == "" {
		code = HTTPCodeConverter(status)
	}
	if message == "" {
		message = defaultErrorMessage(code)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	resp := ErrorResponse{
		Error: ErrorDetails{
			Code:    code,
			Message: message,
		},
	}

	return json.NewEncoder(w).Encode(resp)
}

func defaultErrorMessage(code ErrorCode) string {
	switch code {
	case ErrorCodeInvalidRequest:
		return "invalid request"
	case ErrorCodeUnauthorized:
		return "unauthorized"
	case ErrorCodeForbidden:
		return "forbidden"
	case ErrorCodeNotFound:
		return "not found"
	default:
		return "internal server error"
	}
}
