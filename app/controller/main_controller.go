package controller

import (
	"os"
	"path/filepath"

	"github.com/Web-Engineering-XDU/Project-Backend/app/service"
	agentsystem "github.com/Web-Engineering-XDU/Project-Backend/app/service/agent_system"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/swaggo/files"       // swagger embed files
	"github.com/swaggo/gin-swagger" // gin-swagger middleware
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
	router.GET("/agent-relation/for-edit", service.GetAgentRelationsForEdit)

	ex, err := os.Executable()
    if err != nil {
        panic(err)
    }
	router.Static("/static/rss", filepath.Dir(ex)+"/rss")

	//TODO
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}