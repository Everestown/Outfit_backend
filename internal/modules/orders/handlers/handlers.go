package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Everestown/Outfit_backend/internal/modules/orders/dto"
	"github.com/Everestown/Outfit_backend/internal/modules/orders/service"
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

func (h *Handler) List(c *gin.Context) {
	userID := c.GetUint("user_id")

	orders, err := h.service.GetUserOrders(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch orders"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"orders": orders})
}

func (h *Handler) Get(c *gin.Context) {
	userID := c.GetUint("user_id")
	orderID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order ID"})
		return
	}

	order, err := h.service.GetOrderByID(userID, uint(orderID))
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch order"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"order": order})
}

func (h *Handler) Create(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req dto.CreateOrderRequest
	if !httpx.BindJSON(c, &req) {
		return
	}

	order, err := h.service.CreateOrder(userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create order"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"order": order})
}
