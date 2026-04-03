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
	Address              string
	Env                  string
	BodyLimitBytes       int64
	ReadTimeoutSec       int
	ReadHeaderTimeoutSec int
	WriteTimeoutSec      int
	IdleTimeoutSec       int
	ShutdownTimeoutSec   int
	RateLimitRPS         int
	RateLimitBurst       int
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

	enabledModules := map[string]bool{
		"auth":     viper.GetBool("modules.enabled.auth"),
		"products": viper.GetBool("modules.enabled.products"),
		"cart":     viper.GetBool("modules.enabled.cart"),
		"orders":   viper.GetBool("modules.enabled.orders"),
	}

	return &Config{
		Server: ServerConfig{
			Address:              viper.GetString("server.address"),
			Env:                  viper.GetString("server.env"),
			BodyLimitBytes:       viper.GetInt64("server.body_limit_bytes"),
			ReadTimeoutSec:       viper.GetInt("server.read_timeout_sec"),
			ReadHeaderTimeoutSec: viper.GetInt("server.read_header_timeout_sec"),
			WriteTimeoutSec:      viper.GetInt("server.write_timeout_sec"),
			IdleTimeoutSec:       viper.GetInt("server.idle_timeout_sec"),
			ShutdownTimeoutSec:   viper.GetInt("server.shutdown_timeout_sec"),
			RateLimitRPS:         viper.GetInt("server.rate_limit_rps"),
			RateLimitBurst:       viper.GetInt("server.rate_limit_burst"),
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
