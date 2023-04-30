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

func New(huggoConfig config.Config, setController func(*gin.Engine, *agentsystem.AgentCollection) *gin.Engine) *Huggo {
	agents := agentsystem.NewAgentCollection()
	eventHdl := agentsystem.NewEventHandler()
	ret := &Huggo{
		config: huggoConfig,
		agentSys:  agentsystem.NewAgentSystem(&agents, &eventHdl),
	}
	ret.ginServer = setController(gin.Default(), ret.agentSys.Agents)
	return ret
}

func (huggo *Huggo) Run(port string) {
	if err:=models.InitDB(huggo.config.MySQL) ; err!=nil{
		panic(err)
	}
	go huggo.agentSys.Run()
	huggo.ginServer.Run(port)
}
