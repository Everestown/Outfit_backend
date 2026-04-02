package orders

import (
	"github.com/Everestown/Outfit_backend/internal/core/module"
	"github.com/Everestown/Outfit_backend/internal/models"
	"github.com/Everestown/Outfit_backend/internal/modules/orders/handlers"
	"github.com/Everestown/Outfit_backend/internal/modules/orders/repository"
	"github.com/Everestown/Outfit_backend/internal/modules/orders/service"
	"github.com/Everestown/Outfit_backend/internal/pkg/jwt"
	"github.com/Everestown/Outfit_backend/internal/pkg/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type OrdersModule struct {
	module.BaseModule
	db         *gorm.DB
	jwtManager *jwt.JWTManager
	handler    *handlers.Handler
}

func NewOrdersModule(db *gorm.DB, jwtManager *jwt.JWTManager) module.Module {
	repo := repository.NewRepository(db)
	svc := service.NewService(repo)
	h := handlers.NewHandler(svc)

	return &OrdersModule{
		BaseModule: module.BaseModule{Name: "orders"},
		db:         db,
		jwtManager: jwtManager,
		handler:    h,
	}
}

func (m *OrdersModule) Init() error {
	return m.db.AutoMigrate(
		&models.Order{},
		&models.OrderItem{},
		&models.Payment{},
	)
}

func (m *OrdersModule) RegisterRoutes(router *gin.RouterGroup) {
	ordersGroup := router.Group("/orders")
	ordersGroup.Use(middleware.AuthMiddleware(m.jwtManager))
	{
		ordersGroup.POST("", m.handler.Create)
		ordersGroup.GET("", m.handler.List)
		ordersGroup.GET("/my", m.handler.List)
		ordersGroup.GET("/:id", m.handler.Get)
	}
}

func (m *OrdersModule) Close() error {
	return nil
}
