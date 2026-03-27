package middleware

import (
	"github.com/Everestown/Outfit_backend/internal/config"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CORSMiddleware(cfg *config.CORSConfig) gin.HandlerFunc {
	origins := cfg.AllowedOrigins
	if len(origins) == 0 {
		origins = []string{"http://localhost:3000"}
	}

	return cors.New(cors.Config{
		AllowOrigins:     origins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", RequestIDHeader},
		ExposeHeaders:    []string{"Content-Length", RequestIDHeader},
		AllowCredentials: true,
	})
}
