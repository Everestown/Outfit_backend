package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Everestown/Outfit_backend/internal/modules/cart/dto"
	"github.com/Everestown/Outfit_backend/internal/modules/cart/service"
	"github.com/Everestown/Outfit_backend/internal/pkg/apperrors"
	"github.com/Everestown/Outfit_backend/internal/pkg/httpx"
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
	if !httpx.BindJSON(c, &req) {
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
		switch {
		case errors.Is(err, apperrors.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "item not found"})
		case errors.Is(err, apperrors.ErrForbidden):
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to remove item"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "item removed"})
}
