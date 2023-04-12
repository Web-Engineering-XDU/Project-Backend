package service

import (
	"github.com/gin-gonic/gin"
)

type HuggoConfig struct {
	dataSourceName string
}

type Huggo struct {
	config    HuggoConfig
	ginServer *gin.Engine
	agentSys  *AgentSystem
}

func New(HuggoConfig HuggoConfig) *Huggo {
	agents := newAgentCollection()
	eventHdl := newEventHandler()
	return &Huggo{
		//TODO
		ginServer: gin.Default(),
		agentSys:  newAgentSystem(&agents, &eventHdl),
	}
}

func (huggo *Huggo) Run() {
	// models.InitDB()
	huggo.agentSys.run()
}
