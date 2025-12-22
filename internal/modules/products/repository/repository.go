package repository

import (
	"github.com/Everestown/Outfit_backend/internal/models"
	"gorm.io/gorm"
)

type Repository interface {
	GetAllProducts() ([]models.Product, error)
	GetProductByID(id uint) (*models.Product, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetAllProducts() ([]models.Product, error) {
	var products []models.Product
	err := r.db.
		Preload("Category").
		Preload("Variants.Size").
		Preload("Variants.Color").
		Preload("Images").
		Find(&products).Error
	return products, err
}

func (r *repository) GetProductByID(id uint) (*models.Product, error) {
	var product models.Product
	err := r.db.
		Preload("Category").
		Preload("Variants.Size").
		Preload("Variants.Color").
		Preload("Images").
		First(&product, id).Error
	return &product, err
}
