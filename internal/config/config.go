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
	Env     string
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
	/*enabledModules := make(map[string]bool)
	if rawModules := viper.GetStringMap("modules.enabled"); rawModules != nil {
		for key, value := range rawModules {
			if boolVal, ok := value.(bool); ok {
				enabledModules[key] = boolVal
			}
		}
	}*/
	enabledModules := map[string]bool{
		"auth":     viper.GetBool("modules.enabled.auth"),
		"products": viper.GetBool("modules.enabled.products"),
		"cart":     viper.GetBool("modules.enabled.cart"),
		"orders":   viper.GetBool("modules.enabled.orders"),
	}

	return &Config{
		Server: ServerConfig{
			Address: viper.GetString("server.address"),
			Env:     viper.GetString("server.env"),
		},
		Database: DatabaseConfig{
			URL: viper.GetString("database.url"),
		},
		JWT: JWTConfig{
			Secret: viper.GetString("jwt.secret"),
		},
		Log: LogConfig{
			Level: viper.GetString("log.level"),
		},
		CORS: CORSConfig{
			AllowedOrigins: viper.GetStringSlice("cors.allowed_origins"),
		},
		Modules: ModulesConfig{
			Enabled: enabledModules,
		},
	}
}
