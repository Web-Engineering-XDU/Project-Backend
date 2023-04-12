package service

import (
	"context"
	"time"
)

type agentCollection struct {
	agentMap map[int]*Agent
	eventHdl *eventHandler
}

func newAgentCollection() agentCollection {
	return agentCollection{
		agentMap: make(map[int]*Agent),
	}
}

func (am agentCollection) init() {
	// models.
}

func (agents *agentCollection) NextAgentDo(agentId int, e *Event) {
	agent, ok := agents.agentMap[agentId]
	if !ok {
		//TODO no such agent
		return
	}
	if !agent.AllowInput {
		//target do not allow input
		return
	}
	//agent is not awake. Wake it up
	ctx, cancle := context.WithTimeout(agent.Ctx, 10*time.Minute)
	defer cancle()
	(*agent).Run(ctx, agent ,e)
}
