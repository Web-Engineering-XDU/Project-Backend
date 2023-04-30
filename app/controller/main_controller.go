package controller

import (
	"github.com/Web-Engineering-XDU/Project-Backend/app/service"
	agentsystem "github.com/Web-Engineering-XDU/Project-Backend/app/service/agent_system"
	"github.com/gin-gonic/gin"
)

func SetController(router *gin.Engine, ac *agentsystem.AgentCollection) *gin.Engine {
	router.Use(func(ctx *gin.Context) {
		ctx.Set("agents", ac)
	})
	router.GET("/agent", service.GetAgentList)
	router.PUT("/agent", service.NewAgent)

	return router
}
