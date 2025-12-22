package service

import (
	"github.com/Everestown/Outfit_backend/internal/models"
	"github.com/Everestown/Outfit_backend/internal/modules/cart/repository"
)

type Service interface {
	GetCartByUserID(userID uint) (*models.Cart, error)
	AddItemToCart(userID uint, variantID uint, quantity int) error
	RemoveItemFromCart(userID uint, itemID uint) error
}

type service struct {
	repo repository.Repository
}

func NewService(repo repository.Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetCartByUserID(userID uint) (*models.Cart, error) {
	return s.repo.GetCartByUserID(userID)
}

func (s *service) AddItemToCart(userID uint, variantID uint, quantity int) error {
	return s.repo.AddItemToCart(userID, variantID, quantity)
}

func (s *service) RemoveItemFromCart(userID uint, itemID uint) error {
	return s.repo.RemoveItemFromCart(userID, itemID)
}
