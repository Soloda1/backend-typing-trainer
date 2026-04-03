package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	infraLogger "backend-typing-trainer/internal/infrastructure/logger"
)

func TestRequestLoggerMiddleware(t *testing.T) {
	mw := RequestLoggerMiddleware(infraLogger.New("dev"))
	called := false

	h := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		require.Equal(t, "1", r.URL.Query().Get("page"))
		_, err := w.Write([]byte("ok"))
		require.NoError(t, err)
	}))

	req := httptest.NewRequest(http.MethodGet, "/rooms/list?page=1", http.NoBody)
	req.RemoteAddr = "127.0.0.1:12345"
	req.Header.Set("User-Agent", "middleware-test")
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	require.True(t, called)
	require.Equal(t, http.StatusOK, rr.Code)
	require.Equal(t, "ok", rr.Body.String())
}
