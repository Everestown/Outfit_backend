package module

import "github.com/gin-gonic/gin"

type BaseModule struct {
	Name   string
	Router *gin.RouterGroup
}

func (m *BaseModule) GetName() string {
	return m.Name
}

func (m *BaseModule) RegisterRoutes(router *gin.RouterGroup) {
	m.Router = router
}

func (m *BaseModule) Init() error {
	return nil
}

func (m *BaseModule) Close() error {
	return nil
}
