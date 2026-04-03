package products

import (
	"github.com/Everestown/Outfit_backend/internal/core/module"
	"github.com/Everestown/Outfit_backend/internal/models"
	"github.com/Everestown/Outfit_backend/internal/modules/products/handlers"
	"github.com/Everestown/Outfit_backend/internal/modules/products/repository"
	"github.com/Everestown/Outfit_backend/internal/modules/products/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProductsModule struct {
	module.BaseModule
	db      *gorm.DB
	handler *handlers.Handler
}

func NewProductsModule(db *gorm.DB) module.Module {
	repo := repository.NewRepository(db)
	svc := service.NewService(repo)
	h := handlers.NewHandler(svc)

	return &ProductsModule{
		BaseModule: module.BaseModule{Name: "products"},
		db:         db,
		handler:    h,
	}
}

func (m *ProductsModule) Init() error {
	return m.db.AutoMigrate(
		&models.Category{},
		&models.Product{},
		&models.ProductVariant{},
		&models.Size{},
		&models.Color{},
		&models.ProductImage{},
	)
}

func (m *ProductsModule) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/products", m.handler.List)
	router.GET("/products/:id", m.handler.Get)
	router.GET("/categories", m.handler.ListCategories)
	router.GET("/categories/tree", m.handler.CategoryTree)
}

func (m *ProductsModule) Close() error {
	return nil
}
