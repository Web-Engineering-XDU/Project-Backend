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
	agentRouter := router.Group("/agent")
	agentRouter.PUT("", service.NewAgent)
	agentRouter.DELETE("", service.DeleteAgent)
	agentRouter.POST("",service.UpdateAgent)
	agentRouter.GET("", service.GetAgentList)
	agentRouter.POST("/dry-run", service.DryRun)
	agentRouter.GET("/relationable", service.GetRelationableAgents)

	router.GET("/event", service.GetEventList)

	router.GET("/agent-relation", service.GetAllAgentRelations)
	router.POST("/agent-relation", service.SetAgentRelation)

	//TODO
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}