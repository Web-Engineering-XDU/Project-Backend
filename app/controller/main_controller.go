package controller

import (
	"github.com/Web-Engineering-XDU/Project-Backend/app/service"
	agentsystem "github.com/Web-Engineering-XDU/Project-Backend/app/service/agent_system"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"

	"github.com/swaggo/gin-swagger" // gin-swagger middleware
    "github.com/swaggo/files" // swagger embed files
)

func SetController(router *gin.Engine, ac *agentsystem.AgentCollection) *gin.Engine {
	router.Use(cors.Default())
	router.Use(func(ctx *gin.Context) {
		ctx.Set("agents", ac)
	})
	router.PUT("/agent", service.NewAgent)
	router.DELETE("/agent", service.DeleteAgent)
	router.POST("/agent",service.UpdateAgent)
	router.GET("/agent", service.GetAgentList)

	router.GET("/event", service.GetEventList)

	//TODO
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}