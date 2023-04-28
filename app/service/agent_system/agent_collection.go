package agentsystem

import (
	"context"
	"sync"
	"time"

	"github.com/Web-Engineering-XDU/Project-Backend/app/models"
)

type agentCollection struct {
	agentMap map[int]*Agent
	eventHdl *eventHandler
	ctx      context.Context
}

func NewAgentCollection() agentCollection {
	return agentCollection{
		agentMap: make(map[int]*Agent),
		ctx:      context.Background(),
	}
}

func (ac *agentCollection) init() error {
	agents := models.SelectAgentRuntimeList()
	relations := models.SelectAgentRelationList()

	schedule_agents := make([]int, 0, 10)

	for _, v := range agents {
		ac.agentMap[int(v.ID)] = &Agent{
			AgentInfo: AgentInfo{
				Id:               v.ID,
				Enable:           v.Enable,
				AgentTypeId:      v.TypeId,
				AgentCoreJsonStr: v.PropJsonStr,
				AllowInput:       v.AllowInput,
				AllowOutput:      v.AllowOutput,
				SrcAgentId:       make([]int, 0, 2),
				DstAgentId:       make([]int, 0, 2),
				EventForever:     v.EventForever,
				EventMaxAge:      time.Duration(v.EventMaxAge),
			},
			ac:    ac,
			Ctx:   ac.ctx,
			Mutex: sync.RWMutex{},
		}
		err := ac.agentMap[v.ID].loadCore()
		if err != nil {
			//TODO
			panic(err)
		}
		if v.TypeId == 1 && v.Enable {
			schedule_agents = append(schedule_agents, int(v.ID))
		}
	}

	for _, v := range relations {
		ac.agentMap[v.SrcAgentId].DstAgentId = append(ac.agentMap[v.SrcAgentId].DstAgentId, v.DstAgentId)
		ac.agentMap[v.DstAgentId].SrcAgentId = append(ac.agentMap[v.DstAgentId].SrcAgentId, v.SrcAgentId)
	}

	for _, v := range schedule_agents {
		ac.NextAgentDo(v, nil)
	}

	return nil
}

func (agents *agentCollection) NextAgentDo(agentId int, e *Event) {
	agent, ok := agents.agentMap[agentId]
	if !ok {
		//TODO no such agent
		return
	}
	if !agent.Enable || !agent.AllowInput && agent.AgentTypeId != 1 {
		//target do not allow input
		return
	}
	// ctx, cancle := context.WithTimeout(agent.Ctx, 10*time.Minute)
	// defer cancle()
	// agent.Run(ctx, agent, e)

	agent.Run(context.Background(), agent, e)
}
