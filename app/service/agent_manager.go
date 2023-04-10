package service

import (
	"context"
	"sync"
	"time"
)

type AgentManager struct {
	agentMap map[int]*Agent
}

var am = AgentManager{
	agentMap: make(map[int]*Agent),
}

func GetAgentManager() AgentManager {
	return am
}

func StartAgentManager() {
	StartEventHandler()
	am.init()
}

func (am *AgentManager) init() {
	//TODO
}

var once sync.Once
func NextAgentDo(agentId int, e *Event) {
	agent, ok := am.agentMap[agentId]
	if !ok {
		//no such agent
		return
	}
	if !agent.allowInput {
		//target do not allow input
		return
	}
	agentCore := agent.agentCore
	//agent is not awake. Wake it up
	if agentCore == nil {
		once.Do(func() {
			var err error
			agentCore, err = am.loadAgentCore(agentId)
			if err != nil {
				//TODO cannot wake up exception
			}
		})
	}
	if agentCore == nil {
		//TODO cannot get agent with this Id
		return
	}
	ctx, cancle  := context.WithTimeout(agent.ctx, 10 * time.Minute)
	defer cancle()
	(*agentCore).Run(ctx, agent, e)
}

func (am *AgentManager) loadAgentCore(agentId int) (*AgentCore, error) {
	//TODO
	return nil, nil
}
