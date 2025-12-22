package dto

import "github.com/Everestown/Outfit_backend/internal/models"

type ProductResponse struct {
	ID          uint                    `json:"id"`
	Name        string                  `json:"name"`
	Description string                  `json:"description"`
	BasePrice   float64                 `json:"base_price"`
	Category    models.Category         `json:"category"`
	Variants    []models.ProductVariant `json:"variants"`
	Images      []models.ProductImage   `json:"images"`
}
