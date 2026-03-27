package service

import (
	"errors"

	"github.com/Everestown/Outfit_backend/internal/models"
	"github.com/Everestown/Outfit_backend/internal/modules/products/repository"
)

type CategoryNode struct {
	ID       uint           `json:"id"`
	Name     string         `json:"name"`
	ParentID uint           `json:"parent_id"`
	CatCode  string         `json:"cat_code"`
	Children []CategoryNode `json:"children,omitempty"`
}

type Service interface {
	GetAllProducts() ([]models.Product, error)
	GetProductByID(id uint) (*models.Product, error)
	GetAllCategories() ([]models.Category, error)
	GetCategoryTree() ([]CategoryNode, error)
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

func (s *service) GetAllCategories() ([]models.Category, error) {
	return s.repo.GetAllCategories()
}

func (s *service) GetCategoryTree() ([]CategoryNode, error) {
	categories, err := s.repo.GetAllCategories()
	if err != nil {
		return nil, err
	}

	byParent := make(map[uint][]models.Category, len(categories))
	for _, category := range categories {
		byParent[category.ParentID] = append(byParent[category.ParentID], category)
	}

	return buildCategoryNodes(0, byParent), nil
}

func buildCategoryNodes(parentID uint, byParent map[uint][]models.Category) []CategoryNode {
	nodes := make([]CategoryNode, 0, len(byParent[parentID]))
	for _, category := range byParent[parentID] {
		nodes = append(nodes, CategoryNode{
			ID:       category.ID,
			Name:     category.Name,
			ParentID: category.ParentID,
			CatCode:  category.CatCode,
			Children: buildCategoryNodes(category.ID, byParent),
		})
	}

	return nodes
}
