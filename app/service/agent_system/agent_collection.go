package agentsystem

import (
	"context"
	"sync"
	"time"

	"github.com/Web-Engineering-XDU/Project-Backend/app/models"
)

type AgentCollection struct {
	agentMap     map[int]*Agent
	agentTypeMap map[int]*AgentTypeProp
	eventHdl     *eventHandler
	ctx          context.Context
}

func NewAgentCollection() AgentCollection {
	return AgentCollection{
		agentMap: make(map[int]*Agent),
		ctx:      context.Background(),
	}
}

func (ac *AgentCollection) init() error {

	agentTypes := models.SelectAgentTypeList()
	ac.agentTypeMap = make(map[int]*AgentTypeProp, len(agentTypes))
	for _, v := range agentTypes {
		ac.agentTypeMap[v.ID] = &AgentTypeProp{
			AllowInput:  v.AllowInput,
			AllowOutput: v.AllowOutput,
		}
	}

	agents := models.SelectAgentRuntimeList()
	relations := models.SelectAllAgentRelations()

	schedule_agents := make([]*Agent, 0, 10)

	for _, v := range agents {
		ac.agentMap[v.ID] = &Agent{
			AgentInfo: AgentInfo{
				ID:               v.ID,
				Enable:           v.Enable,
				TypeId:           v.TypeId,
				AgentCoreJsonStr: v.PropJsonStr,
				SrcAgentId:       make([]int, 0, 2),
				DstAgentId:       make([]int, 0, 2),
				EventForever:     v.EventForever,
				EventMaxAge:      time.Duration(v.EventMaxAge),
			},
			ac:    ac,
			Ctx:   ac.ctx,
			Mutex: sync.RWMutex{},
		}
		err := ac.agentMap[v.ID].LoadCore()
		if err != nil {
			//TODO
			panic(err)
		}
		if v.TypeId == ScheduleAgentId && v.Enable {
			schedule_agents = append(schedule_agents, ac.agentMap[v.ID])
		}
	}

	for _, v := range relations {
		//TODO
		_, ok := ac.agentMap[v.SrcAgentId]
		if ok {
			ac.agentMap[v.SrcAgentId].DstAgentId = append(ac.agentMap[v.SrcAgentId].DstAgentId, v.DstAgentId)
		}
		_, ok = ac.agentMap[v.DstAgentId]
		if ok {
			ac.agentMap[v.DstAgentId].SrcAgentId = append(ac.agentMap[v.DstAgentId].SrcAgentId, v.SrcAgentId)
		}
	}

	for _, v := range schedule_agents {
		go v.Run(ac.ctx, v, nil, ac.eventHdl.PushEvent)
	}

	return nil
}

func (agents *AgentCollection) NextAgentDo(agentId int, e *Event) {
	agent, ok := agents.agentMap[agentId]
	if !ok {
		//TODO no such agent
		return
	}
	agentTypeInfo, ok := agents.agentTypeMap[agent.TypeId]
	if !ok {
		//TODO no such type
		return
	}
	if !agent.Enable || !agentTypeInfo.AllowInput {
		//target do not allow input
		return
	}
	// ctx, cancle := context.WithTimeout(agent.Ctx, 10*time.Minute)
	// defer cancle()
	// agent.Run(ctx, agent, e)

	agent.Run(agent.Ctx, agent, e, agents.eventHdl.PushEvent)
}
