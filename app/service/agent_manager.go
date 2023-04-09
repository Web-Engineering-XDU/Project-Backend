package service

import "github.com/Web-Engineering-XDU/Project-Backend/app/service/interfaces"

type AgentMata struct {
	id int
	name string	

	agent *interfaces.Agent
}

type AgentManager struct {
	agentMap map[int]AgentMata
}

var am = AgentManager{
	agentMap: make(map[int]AgentMata),
}

func GetAgentManager() *AgentManager {
	return &am
}

func StartAgentManager(){
	am.init()
}

func (am *AgentManager) init() {

}