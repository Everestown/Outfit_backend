package handlers

import (
	"net/http"
	"strconv"

	"github.com/Everestown/Outfit_backend/internal/modules/cart/dto"
	"github.com/Everestown/Outfit_backend/internal/modules/cart/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service service.Service
}

func NewHandler(service service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetCart(c *gin.Context) {
	userID := c.GetUint("user_id")

	cart, err := h.service.GetCartByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch cart"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"cart": cart})
}

func (h *Handler) AddItem(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req dto.AddItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.AddItemToCart(userID, req.VariantID, req.Quantity); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add item"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "item added"})
}

func (h *Handler) RemoveItem(c *gin.Context) {
	userID := c.GetUint("user_id")

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid item ID"})
		return
	}

	if err := h.service.RemoveItemFromCart(userID, uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to remove item"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "item removed"})
}
