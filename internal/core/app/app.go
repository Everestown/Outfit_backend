package app

import (
	"github.com/Everestown/Outfit_backend/internal/modules/auth"
	"github.com/Everestown/Outfit_backend/internal/modules/cart"
	"github.com/Everestown/Outfit_backend/internal/modules/orders"
	"github.com/Everestown/Outfit_backend/internal/modules/products"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/Everestown/Outfit_backend/internal/config"
	"github.com/Everestown/Outfit_backend/internal/core/module"
	"github.com/Everestown/Outfit_backend/internal/logger"
	"github.com/Everestown/Outfit_backend/internal/pkg/database"
	"github.com/Everestown/Outfit_backend/internal/pkg/jwt"
	"github.com/Everestown/Outfit_backend/internal/pkg/middleware"
	"github.com/Everestown/Outfit_backend/internal/pkg/swagger"
)

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
		l.Fatal("Failed to connect to database", logger.Error(err))
	}

	jwtManager := jwt.NewJWTManager(cfg.JWT.Secret, db)

	router := gin.Default()

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
	a.router.Use(gin.Recovery())
	a.router.Use(middleware.CORSMiddleware(&a.config.CORS))

	if a.config.Server.Env == "development" {
		swagger.SetupSwagger(a.router)
	}
}

func (a *App) RegisterCoreModules() {
	// Регистрация стандартных модулей из конфига
	if a.config.Modules.Enabled["auth"] {
		err := a.RegisterModule(auth.NewAuthModule(a.db, a.jwt))
		if err != nil {
			return
		}
	}
	if a.config.Modules.Enabled["products"] {
		err := a.RegisterModule(products.NewProductsModule(a.db))
		if err != nil {
			return
		}
	}
	if a.config.Modules.Enabled["cart"] {
		err := a.RegisterModule(cart.NewCartModule(a.db))
		if err != nil {
			return
		}
	}
	if a.config.Modules.Enabled["orders"] {
		err := a.RegisterModule(orders.NewOrdersModule(a.db))
		if err != nil {
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
	api := a.router.Group("/api")

	for _, m := range a.registry.GetAllModules() {
		m.RegisterRoutes(api.Group("/" + m.GetName()))
	}
}

func (a *App) Run() error {
	a.SetupRouter()
	return a.router.Run(a.config.Server.Address)
}

func (a *App) Close() {
	a.registry.CloseAll()
}
