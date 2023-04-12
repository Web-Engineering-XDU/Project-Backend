package service

import (
	"github.com/Web-Engineering-XDU/Project-Backend/app/models"
	"github.com/Web-Engineering-XDU/Project-Backend/config"
	"github.com/gin-gonic/gin"
)

type Huggo struct {
	config    config.Config
	ginServer *gin.Engine
	agentSys  *AgentSystem
}

func New(huggoConfig config.Config) *Huggo {
	agents := newAgentCollection()
	eventHdl := newEventHandler()
	return &Huggo{
		config: huggoConfig,
		//TODO
		ginServer: gin.Default(),
		agentSys:  newAgentSystem(&agents, &eventHdl),
	}
}

func (huggo *Huggo) Run() {
	models.InitDB(huggo.config.MySQL)
	huggo.agentSys.run()
}
