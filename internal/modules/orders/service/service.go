package service

import (
	"errors"

	"github.com/Everestown/Outfit_backend/internal/models"
	"github.com/Everestown/Outfit_backend/internal/modules/orders/dto"
	"github.com/Everestown/Outfit_backend/internal/modules/orders/repository"
)

type Service interface {
	GetUserOrders(userID uint) ([]models.Order, error)
	GetOrderByID(userID uint, orderID uint) (*models.Order, error)
	CreateOrder(userID uint, req dto.CreateOrderRequest) (*models.Order, error)
}

type service struct {
	repo repository.Repository
}

func NewService(repo repository.Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetUserOrders(userID uint) ([]models.Order, error) {
	return s.repo.GetUserOrders(userID)
}

func (s *service) GetOrderByID(userID uint, orderID uint) (*models.Order, error) {
	order, err := s.repo.GetOrderByID(userID, orderID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, errors.New("order not found")
	}
	return order, nil
}

func (s *service) CreateOrder(userID uint, req dto.CreateOrderRequest) (*models.Order, error) {
	return s.repo.CreateOrder(userID, req)
}
