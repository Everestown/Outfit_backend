package handlers

import (
	"net/http"

	"github.com/Everestown/Outfit_backend/internal/modules/auth/dto"
	"github.com/Everestown/Outfit_backend/internal/modules/auth/service"
	"github.com/Everestown/Outfit_backend/internal/pkg/httpx"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service service.Service
}

func NewHandler(s service.Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if !httpx.BindJSON(c, &req) {
		return
	}

	ctx := dto.SessionContext{
		IP:         getClientIP(c),
		DeviceInfo: c.GetHeader("User-Agent"),
	}

	tokens, err := h.service.Register(req, ctx)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, tokens)
}

func (h *Handler) Login(c *gin.Context) {
	var req dto.LoginRequest

	if !httpx.BindJSON(c, &req) {
		return
	}

	ctx := dto.SessionContext{
		IP:         getClientIP(c),
		DeviceInfo: c.GetHeader("User-Agent"),
	}

	tokens, err := h.service.Login(req, ctx)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tokens)
}

func (h *Handler) Refresh(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	if !httpx.BindJSON(c, &req) {
		return
	}

	ctx := dto.SessionContext{
		IP:         getClientIP(c),
		DeviceInfo: c.GetHeader("User-Agent"),
	}

	tokens, err := h.service.RefreshToken(req.RefreshToken, ctx)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tokens)
}

func (h *Handler) Logout(c *gin.Context) {
	userID := c.GetUint("user_id")

	if err := h.service.Logout(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "logged out"})
}

func (h *Handler) Profile(c *gin.Context) {
	userID := c.GetUint("user_id")

	user, err := h.service.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

func getClientIP(c *gin.Context) string {
	ip := c.ClientIP()
	if ip == "::1" {
		return "127.0.0.1"
	}
	return ip
}
