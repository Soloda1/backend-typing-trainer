package utils

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHTTPCodeConverter(t *testing.T) {
	require.Equal(t, ErrorCodeInvalidRequest, HTTPCodeConverter(http.StatusBadRequest))
	require.Equal(t, ErrorCodeUnauthorized, HTTPCodeConverter(http.StatusUnauthorized))
	require.Equal(t, ErrorCodeForbidden, HTTPCodeConverter(http.StatusForbidden))
	require.Equal(t, ErrorCodeNotFound, HTTPCodeConverter(http.StatusNotFound))
	require.Equal(t, ErrorCodeInternalError, HTTPCodeConverter(http.StatusConflict))
}

func TestWriteJSON(t *testing.T) {
	rr := httptest.NewRecorder()

	err := WriteJSON(rr, http.StatusCreated, map[string]string{"status": "ok"})
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, rr.Code)
	require.Equal(t, "application/json", rr.Header().Get("Content-Type"))

	var body map[string]string
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &body))
	require.Equal(t, "ok", body["status"])
}

func TestWriteError_UsesProvidedCodeAndMessage(t *testing.T) {
	rr := httptest.NewRecorder()

	err := WriteError(rr, http.StatusConflict, ErrorCodeInvalidRequest, "custom error")
	require.NoError(t, err)
	require.Equal(t, http.StatusConflict, rr.Code)

	var body ErrorResponse
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &body))
	require.Equal(t, ErrorCodeInvalidRequest, body.Error.Code)
	require.Equal(t, "custom error", body.Error.Message)
}

func TestWriteError_DefaultsByStatusAndCode(t *testing.T) {
	rr := httptest.NewRecorder()

	err := WriteError(rr, http.StatusBadRequest, "", "")
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, rr.Code)

	var body ErrorResponse
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &body))
	require.Equal(t, ErrorCodeInvalidRequest, body.Error.Code)
	require.Equal(t, "invalid request", body.Error.Message)
}

func TestDefaultErrorMessage(t *testing.T) {
	require.Equal(t, "invalid request", defaultErrorMessage(ErrorCodeInvalidRequest))
	require.Equal(t, "unauthorized", defaultErrorMessage(ErrorCodeUnauthorized))
	require.Equal(t, "forbidden", defaultErrorMessage(ErrorCodeForbidden))
	require.Equal(t, "internal server error", defaultErrorMessage(ErrorCodeInternalError))
}
