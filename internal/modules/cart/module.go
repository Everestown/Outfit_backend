package cart

import (
	"github.com/Everestown/Outfit_backend/internal/core/module"
	"github.com/Everestown/Outfit_backend/internal/models"
	"github.com/Everestown/Outfit_backend/internal/modules/cart/handlers"
	"github.com/Everestown/Outfit_backend/internal/modules/cart/repository"
	"github.com/Everestown/Outfit_backend/internal/modules/cart/service"
	"github.com/Everestown/Outfit_backend/internal/pkg/jwt"
	"github.com/Everestown/Outfit_backend/internal/pkg/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CartModule struct {
	module.BaseModule
	db         *gorm.DB
	jwtManager *jwt.JWTManager
	handler    *handlers.Handler
}

func NewCartModule(db *gorm.DB, jwtManager *jwt.JWTManager) module.Module {
	repo := repository.NewRepository(db)
	svc := service.NewService(repo)
	h := handlers.NewHandler(svc)

	return &CartModule{
		BaseModule: module.BaseModule{Name: "cart"},
		db:         db,
		jwtManager: jwtManager,
		handler:    h,
	}
}

func (m *CartModule) Init() error {
	return m.db.AutoMigrate(&models.Cart{}, &models.CartItem{})
}

func (m *CartModule) RegisterRoutes(router *gin.RouterGroup) {
	cartGroup := router.Group("/cart")
	cartGroup.Use(middleware.AuthMiddleware(m.jwtManager))
	{
		cartGroup.GET("", m.handler.GetCart)
		cartGroup.POST("/items", m.handler.AddItem)
		cartGroup.DELETE("/items/:id", m.handler.RemoveItem)
	}
}

func (m *CartModule) Close() error {
	return nil
}
