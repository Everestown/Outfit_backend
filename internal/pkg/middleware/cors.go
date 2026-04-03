package middleware

import (
	"github.com/Everestown/Outfit_backend/internal/config"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CORSMiddleware(cfg *config.CORSConfig) gin.HandlerFunc {
	origins := cfg.AllowedOrigins
	if len(origins) == 0 && gin.Mode() != gin.ReleaseMode {
		origins = []string{"http://localhost:3000"}
	}

	corsCfg := cors.Config{
		AllowOrigins:     origins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", RequestIDHeader},
		ExposeHeaders:    []string{"Content-Length", RequestIDHeader},
		AllowCredentials: true,
	}

	if len(origins) == 0 {
		corsCfg.AllowOriginFunc = func(string) bool { return false }
	}

	return cors.New(corsCfg)
}
