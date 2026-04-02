package app

import (
	"net/http"

	"github.com/Everestown/Outfit_backend/internal/modules/auth"
	"github.com/Everestown/Outfit_backend/internal/modules/cart"
	"github.com/Everestown/Outfit_backend/internal/modules/orders"
	"github.com/Everestown/Outfit_backend/internal/modules/products"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/Everestown/Outfit_backend/internal/config"
	"github.com/Everestown/Outfit_backend/internal/core/module"
	"github.com/Everestown/Outfit_backend/internal/logger"
	"github.com/Everestown/Outfit_backend/internal/pkg/database"
	"github.com/Everestown/Outfit_backend/internal/pkg/jwt"
	"github.com/Everestown/Outfit_backend/internal/pkg/middleware"
	"github.com/Everestown/Outfit_backend/internal/pkg/swagger"
)

const defaultBodyLimitBytes int64 = 2 << 20 // 2 MiB

type App struct {
	config   *config.Config
	router   *gin.Engine
	db       *gorm.DB
	jwt      *jwt.JWTManager
	registry *ModuleRegistry
	logger   *logger.Logger
}

func NewApp(cfg *config.Config) *App {
	l := logger.New(cfg.Log.Level)

	db, err := database.NewPostgresDB(&cfg.Database)
	if err != nil {
		l.Fatal("Failed to connect to database", logger.Err(err))
	}

	jwtManager := jwt.NewJWTManager(cfg.JWT.Secret, db)

	if cfg.Server.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	app := &App{
		config:   cfg,
		router:   router,
		db:       db,
		jwt:      jwtManager,
		registry: NewModuleRegistry(),
		logger:   l,
	}

	app.setupMiddleware()

	return app
}

func (a *App) setupMiddleware() {
	_ = a.router.SetTrustedProxies(nil)

	a.router.Use(middleware.RequestIDMiddleware())
	a.router.Use(middleware.SecurityHeadersMiddleware())
	a.router.Use(middleware.BodyLimitMiddleware(defaultBodyLimitBytes))
	a.router.Use(gin.Recovery())
	a.router.Use(middleware.CORSMiddleware(&a.config.CORS))

	a.router.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	if a.config.Server.Env == "development" {
		swagger.SetupSwagger(a.router)
	}
}

func (a *App) RegisterCoreModules() {
	a.logger.Info("Modules enabled", zap.Any("enabled", a.config.Modules.Enabled))

	if a.config.Modules.Enabled["auth"] {
		err := a.RegisterModule(auth.NewAuthModule(a.db, a.jwt))
		if err != nil {
			a.logger.Error("Failed to register auth module", logger.Err(err))
			return
		}
	}
	if a.config.Modules.Enabled["products"] {
		err := a.RegisterModule(products.NewProductsModule(a.db))
		if err != nil {
			a.logger.Error("Failed to register products module", logger.Err(err))
			return
		}
	}
	if a.config.Modules.Enabled["cart"] {
		err := a.RegisterModule(cart.NewCartModule(a.db, a.jwt))
		if err != nil {
			a.logger.Error("Failed to register cart module", logger.Err(err))
			return
		}
	}
	if a.config.Modules.Enabled["orders"] {
		err := a.RegisterModule(orders.NewOrdersModule(a.db, a.jwt))
		if err != nil {
			a.logger.Error("Failed to register orders module", logger.Err(err))
			return
		}
	}
}

func (a *App) RegisterModule(m module.Module) error {
	return a.registry.RegisterModule(m)
}

func (a *App) GetDB() *gorm.DB {
	return a.db
}

func (a *App) GetJWT() *jwt.JWTManager {
	return a.jwt
}

func (a *App) SetupRouter() {
	modules := a.registry.GetAllModules()
	a.logger.Info("Active modules count", zap.Int("count", len(modules)))

	api := a.router.Group("/api")

	for _, m := range a.registry.GetAllModules() {
		m.RegisterRoutes(api)
	}
}

func (a *App) Run() error {
	a.SetupRouter()
	return a.router.Run(a.config.Server.Address)
}

func (a *App) Close() {
	a.registry.CloseAll()
}
