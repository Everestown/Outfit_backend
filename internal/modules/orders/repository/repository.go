package repository

import (
	"time"

	"github.com/Everestown/Outfit_backend/internal/models"
	"github.com/Everestown/Outfit_backend/internal/modules/orders/dto"
	"gorm.io/gorm"
)

type Repository interface {
	GetUserOrders(userID uint) ([]models.Order, error)
	CreateOrder(userID uint, req dto.CreateOrderRequest) (*models.Order, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetUserOrders(userID uint) ([]models.Order, error) {
	var orders []models.Order
	err := r.db.
		Preload("Items.Variant").
		Preload("Payment").
		Where("user_id = ?", userID).
		Find(&orders).Error
	return orders, err
}

func (r *repository) CreateOrder(userID uint, req dto.CreateOrderRequest) (*models.Order, error) {
	order := &models.Order{
		UserID:      userID,
		TotalAmount: req.TotalAmount,
		Status:      "pending",
		ShippingAddress: models.ShippingAddress{
			Postcode:   req.Postcode,
			Country:    req.Country,
			City:       req.City,
			Street:     req.Street,
			House:      req.House,
			Apartment:  req.Apartment,
			Phone:      req.Phone,
			Surname:    req.Surname,
			Name:       req.Name,
			Patronymic: req.Patronymic,
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := r.db.Create(order).Error; err != nil {
		return nil, err
	}

	return order, nil
}
