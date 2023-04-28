package service

import (
	"github.com/Web-Engineering-XDU/Project-Backend/app/models"
	agentsystem "github.com/Web-Engineering-XDU/Project-Backend/app/service/agent_system"
	"github.com/Web-Engineering-XDU/Project-Backend/config"
	"github.com/gin-gonic/gin"
)

type Huggo struct {
	config    config.Config
	ginServer *gin.Engine
	agentSys  *agentsystem.AgentSystem
}

func New(huggoConfig config.Config) *Huggo {
	agents := agentsystem.NewAgentCollection()
	eventHdl := agentsystem.NewEventHandler()
	return &Huggo{
		config: huggoConfig,
		//TODO
		ginServer: gin.Default(),
		agentSys:  agentsystem.NewAgentSystem(&agents, &eventHdl),
	}
}

func (huggo *Huggo) Run() {
	if err:=models.InitDB(huggo.config.MySQL) ; err!=nil{
		panic(err)
	}
	huggo.agentSys.Run()
}
