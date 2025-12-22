package module

import "github.com/gin-gonic/gin"

type Module interface {
	GetName() string
	RegisterRoutes(router *gin.RouterGroup)
	Init() error
	Close() error
}
