package orders

import (
	"github.com/Everestown/Outfit_backend/internal/core/module"
	"github.com/Everestown/Outfit_backend/internal/models"
	"github.com/Everestown/Outfit_backend/internal/modules/orders/handlers"
	"github.com/Everestown/Outfit_backend/internal/modules/orders/repository"
	"github.com/Everestown/Outfit_backend/internal/modules/orders/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type OrdersModule struct {
	module.BaseModule
	db      *gorm.DB
	handler *handlers.Handler
}

func NewOrdersModule(db *gorm.DB) module.Module {
	repo := repository.NewRepository(db)
	svc := service.NewService(repo)
	h := handlers.NewHandler(svc)

	return &OrdersModule{
		BaseModule: module.BaseModule{Name: "orders"},
		db:         db,
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
	{
		ordersGroup.GET("", m.handler.List)
		ordersGroup.POST("", m.handler.Create)
	}
}

func (m *OrdersModule) Close() error {
	return nil
}
