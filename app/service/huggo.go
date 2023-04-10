package service

import (
	"github.com/gin-gonic/gin"
)

type HuggoConfig struct {
}

type Huggo struct {
	config    HuggoConfig
	ginServer *gin.Engine
	agentSys  *AgentSystem
}

func New() *Huggo {
	agents := newAgentCollection()
	eventHdl := newEventHandler()
	return &Huggo{
		//TODO
		ginServer: gin.Default(),
		agentSys:  newAgentSystem(&agents, &eventHdl),
	}
}

func (*Huggo) Run() {

}
