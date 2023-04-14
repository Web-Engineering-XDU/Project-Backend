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
	agents, err := models.DB().GetAgentRuntimeInfoList(ac.ctx)
	if err != nil {
		return err
	}
	relations, err := models.DB().GetAgentRelationList(ac.ctx)
	if err != nil {
		return err
	}

	schedule_agents := make([]int, 0, 10)

	for _, v := range agents {
		ac.agentMap[int(v.ID)] = &Agent{
			AgentInfo: AgentInfo{
				Id:               int(v.ID),
				Enable:           v.Enable,
				AgentTypeId:      int(v.TypeID),
				AgentCoreJsonStr: v.PropJsonStr,
				AllowInput:       v.AllowInput,
				AllowOutput:      v.AllowOutput,
				SrcAgentId:       make([]int, 0, 2),
				DstAgentId:       make([]int, 0, 2),
				EventMaxAge:      time.Duration(v.EventMaxAge),
			},
			ac:    ac,
			Ctx:   ac.ctx,
			Mutex: sync.RWMutex{},
		}
		err = ac.agentMap[int(v.ID)].loadCore()
		if err != nil {
			//TODO
			panic(err)
		}
		if v.TypeID == 1 && v.Enable {
			schedule_agents = append(schedule_agents, int(v.ID))
		}
	}

	for _, v := range relations {
		ac.agentMap[int(v.SrcAgentID)].DstAgentId = append(ac.agentMap[int(v.SrcAgentID)].DstAgentId, int(v.DstAgentID))
		ac.agentMap[int(v.DstAgentID)].SrcAgentId = append(ac.agentMap[int(v.DstAgentID)].SrcAgentId, int(v.SrcAgentID))
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
	if !agent.AllowInput && agent.AgentTypeId != 1 {
		//target do not allow input
		return
	}
	// ctx, cancle := context.WithTimeout(agent.Ctx, 10*time.Minute)
	// defer cancle()
	// agent.Run(ctx, agent, e)

	agent.Run(context.Background(), agent , e)
}
