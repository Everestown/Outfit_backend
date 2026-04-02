package auth

import (
	"github.com/Everestown/Outfit_backend/internal/core/module"
	"github.com/Everestown/Outfit_backend/internal/models"
	"github.com/Everestown/Outfit_backend/internal/modules/auth/handlers"
	"github.com/Everestown/Outfit_backend/internal/modules/auth/repository"
	"github.com/Everestown/Outfit_backend/internal/modules/auth/service"
	"github.com/Everestown/Outfit_backend/internal/pkg/jwt"
	"github.com/Everestown/Outfit_backend/internal/pkg/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthModule struct {
	module.BaseModule
	db         *gorm.DB
	jwtManager *jwt.JWTManager
	handler    *handlers.Handler
}

func NewAuthModule(db *gorm.DB, jwtManager *jwt.JWTManager) module.Module {
	repo := repository.NewRepository(db)
	svc := service.NewService(repo, jwtManager)
	h := handlers.NewHandler(svc)

	return &AuthModule{
		BaseModule: module.BaseModule{Name: "auth"},
		db:         db,
		jwtManager: jwtManager,
		handler:    h,
	}
}

func (m *AuthModule) Init() error {
	return m.db.AutoMigrate(&models.User{}, &models.UserSession{})
}

func (m *AuthModule) RegisterRoutes(router *gin.RouterGroup) {
	authGroup := router.Group("/auth")
	authGroup.POST("/register", m.handler.Register)
	authGroup.POST("/login", m.handler.Login)
	authGroup.POST("/refresh", m.handler.Refresh)

	protected := authGroup.Group("")
	protected.Use(middleware.AuthMiddleware(m.jwtManager))
	{
		protected.POST("/logout", m.handler.Logout)
		protected.GET("/profile", m.handler.Profile)
	}
}

func (m *AuthModule) Close() error {
	return nil
}
