package service

import (
	"errors"
	"github.com/Everestown/Outfit_backend/internal/modules/products/repository"

	"github.com/Everestown/Outfit_backend/internal/models"
)

type Service interface {
	GetAllProducts() ([]models.Product, error)
	GetProductByID(id uint) (*models.Product, error)
}

type service struct {
	repo repository.Repository
}

func NewService(repo repository.Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetAllProducts() ([]models.Product, error) {
	return s.repo.GetAllProducts()
}

func (s *service) GetProductByID(id uint) (*models.Product, error) {
	product, err := s.repo.GetProductByID(id)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, errors.New("product not found")
	}
	return product, nil
}
