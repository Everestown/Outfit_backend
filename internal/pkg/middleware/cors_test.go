package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Everestown/Outfit_backend/internal/config"
	"github.com/gin-gonic/gin"
)

func TestCORSMiddleware_DefaultLocalhostForNonRelease(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(CORSMiddleware(&config.CORSConfig{}))
	r.OPTIONS("/ping", func(c *gin.Context) { c.Status(http.StatusNoContent) })

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodOptions, "/ping", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	req.Header.Set("Access-Control-Request-Method", "GET")
	r.ServeHTTP(w, req)

	if got := w.Header().Get("Access-Control-Allow-Origin"); got != "http://localhost:3000" {
		t.Fatalf("expected localhost origin allowed, got %q", got)
	}
}

func TestCORSMiddleware_ReleaseModeNoFallbackOrigins(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	t.Cleanup(func() { gin.SetMode(gin.TestMode) })

	r := gin.New()
	r.Use(CORSMiddleware(&config.CORSConfig{}))
	r.OPTIONS("/ping", func(c *gin.Context) { c.Status(http.StatusNoContent) })

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodOptions, "/ping", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	req.Header.Set("Access-Control-Request-Method", "GET")
	r.ServeHTTP(w, req)

	if got := w.Header().Get("Access-Control-Allow-Origin"); got != "" {
		t.Fatalf("expected no allowed origin in release fallback mode, got %q", got)
	}
}
