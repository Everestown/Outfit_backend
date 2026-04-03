package httpx

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

type bindReq struct {
	Name string `json:"name" binding:"required"`
}

func TestBindJSON_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	body, _ := json.Marshal(bindReq{Name: "ok"})
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req

	var payload bindReq
	if ok := BindJSON(c, &payload); !ok {
		t.Fatalf("expected bind success")
	}
	if payload.Name != "ok" {
		t.Fatalf("unexpected payload: %q", payload.Name)
	}
}

func TestBindJSON_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(`{"name":`))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req

	var payload bindReq
	if ok := BindJSON(c, &payload); ok {
		t.Fatalf("expected bind failure")
	}
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestBindJSON_BodyTooLarge(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	bigBody := bytes.NewBufferString(`{"name":"toolong"}`)
	req := httptest.NewRequest(http.MethodPost, "/", bigBody)
	req.Header.Set("Content-Type", "application/json")
	req.Body = http.MaxBytesReader(w, req.Body, 5)
	c.Request = req

	var payload bindReq
	if ok := BindJSON(c, &payload); ok {
		t.Fatalf("expected bind failure")
	}
	if w.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("expected 413, got %d", w.Code)
	}
}
