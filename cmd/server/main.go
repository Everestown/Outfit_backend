package main

import (
	"github.com/Everestown/Outfit_backend/internal/config"
	"github.com/Everestown/Outfit_backend/internal/core/app"
	"github.com/Everestown/Outfit_backend/internal/logger"
)

func main() {
	cfg := config.Load()

	application := app.NewApp(cfg)

	// Регистрация основных модулей (из конфига)
	application.RegisterCoreModules()

	// 🔌 Plug-and-play: добавление нового модуля
	// application.RegisterModule(payments.NewPaymentsModule(application.GetDB()))

	// Запуск приложения
	if err := application.Run(); err != nil {
		logger.Fatal("Failed to run app", logger.Err(err))
	}
}
