package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRateLimitMiddleware_EnforcesLimit(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(RateLimitMiddleware(1, 1))
	r.GET("/ping", func(c *gin.Context) { c.Status(http.StatusOK) })

	w1 := httptest.NewRecorder()
	req1 := httptest.NewRequest(http.MethodGet, "/ping", nil)
	req1.RemoteAddr = "127.0.0.1:12345"
	r.ServeHTTP(w1, req1)
	if w1.Code != http.StatusOK {
		t.Fatalf("expected first request 200, got %d", w1.Code)
	}

	w2 := httptest.NewRecorder()
	req2 := httptest.NewRequest(http.MethodGet, "/ping", nil)
	req2.RemoteAddr = "127.0.0.1:12345"
	r.ServeHTTP(w2, req2)
	if w2.Code != http.StatusTooManyRequests {
		t.Fatalf("expected second request 429, got %d", w2.Code)
	}
}
