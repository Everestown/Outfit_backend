package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Log      LogConfig
	CORS     CORSConfig
	Modules  ModulesConfig
}

type ServerConfig struct {
	Address string
}
type DatabaseConfig struct {
	URL string
}
type JWTConfig struct {
	Secret string
}
type LogConfig struct {
	Level string
}
type CORSConfig struct {
	AllowedOrigins []string
}
type ModulesConfig struct {
	Enabled map[string]bool
}

func Load() *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		// Используем env vars
	}

	// Преобразование map[string]any в map[string]bool
	enabledModules := make(map[string]bool)
	if rawModules := viper.GetStringMap("modules.enabled"); rawModules != nil {
		for key, value := range rawModules {
			if boolVal, ok := value.(bool); ok {
				enabledModules[key] = boolVal
			}
		}
	}

	return &Config{
		Server: ServerConfig{
			Address: viper.GetString("SERVER_ADDRESS"),
		},
		Database: DatabaseConfig{
			URL: viper.GetString("DATABASE_URL"),
		},
		JWT: JWTConfig{
			Secret: viper.GetString("JWT_SECRET"),
		},
		Log: LogConfig{
			Level: viper.GetString("LOG_LEVEL"),
		},
		CORS: CORSConfig{
			AllowedOrigins: viper.GetStringSlice("CORS_ALLOWED_ORIGINS"),
		},
		Modules: ModulesConfig{
			Enabled: enabledModules,
		},
	}
}
