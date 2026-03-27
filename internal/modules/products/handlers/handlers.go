package handlers

import (
	"net/http"
	"strconv"

	"github.com/Everestown/Outfit_backend/internal/modules/products/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service service.Service
}

func NewHandler(service service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) List(c *gin.Context) {
	products, err := h.service.GetAllProducts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch products"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"products": products})
}

func (h *Handler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product ID"})
		return
	}

	product, err := h.service.GetProductByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"product": product})
}

func (h *Handler) ListCategories(c *gin.Context) {
	categories, err := h.service.GetAllCategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch categories"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"categories": categories})
}

func (h *Handler) CategoryTree(c *gin.Context) {
	categories, err := h.service.GetCategoryTree()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to build category tree"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"categories": categories})
}
