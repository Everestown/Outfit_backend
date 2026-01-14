package handlers

import (
	"fmt"
	"net/http"

	"github.com/Everestown/Outfit_backend/internal/modules/auth/dto"
	"github.com/Everestown/Outfit_backend/internal/modules/auth/service"
	"github.com/Everestown/Outfit_backend/internal/pkg/jwt"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service    service.Service
	jwtManager *jwt.JWTManager
}

func NewHandler(s service.Service, jwtManager *jwt.JWTManager) *Handler {
	return &Handler{
		service:    s,
		jwtManager: jwtManager,
	}
}

func (h *Handler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Printf("REQ: %+v\n", req)

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
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if err := h.service.Logout(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "logged out"})
}

func (h *Handler) Profile(c *gin.Context) {
	userID := c.GetUint("user_id")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	user, err := h.service.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

func getClientIP(c *gin.Context) string {
	ip := c.ClientIP()

	// localhost "::1"
	if ip == "::1" {
		return "127.0.0.1"
	}

	return ip
}
