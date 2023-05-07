package agentsystem

import (
	"errors"
	"sync"
	"time"

	"github.com/Web-Engineering-XDU/Project-Backend/app/models"
)

func (ac *AgentCollection) AddAgent(a *models.Agent) error {
	_, ok := ac.agentMap[a.ID]
	if ok {
		//Already running
		return errors.New("agent with this id already exists")
	}
	_, ok = ac.agentTypeMap[a.TypeId]
	if !ok {
		return errors.New("agent type with this id does not exist")
	}

	agent := &Agent{
		AgentInfo: AgentInfo{
			ID:               a.ID,
			Enable:           a.Enable,
			TypeId:           a.TypeId,
			AgentCoreJsonStr: a.PropJsonStr,
			SrcAgentId:       make([]int, 0, 2),
			DstAgentId:       make([]int, 0, 2),
			EventForever:     a.EventForever,
			EventMaxAge:      time.Duration(a.EventMaxAge),
		},
		ac:    ac,
		Ctx:   ac.ctx,
		Mutex: sync.RWMutex{},
	}
	err := agent.loadCore()
	if err != nil {
		return err
	}

	err = models.InsertAgent(a)
	if err != nil {
		return err
	}

	agent.ID = a.ID
	ac.agentMap[agent.ID] = agent
	if agent.Enable && agent.TypeId == ScheduleAgentId {
		go agent.Run(agent.Ctx, agent, nil, ac.eventHdl.PushEvent)
	}
	return nil
}

func (ac *AgentCollection) DeleteAgent(id int) (bool, error) {
	agent, ok := ac.agentMap[id]
	if !ok {
		return false, nil
	}

	ok, err := models.DeleteAgentAndRelationsAbout(id)
	if err != nil {
		return ok, err
	}

	agent.Stop()
	delete(ac.agentMap, id)
	return true, nil
}

func (ac *AgentCollection) UpdateAgent(a models.Agent) error {
	var err error

	agent, ok := ac.agentMap[a.ID]
	if !ok {
		return errors.New("agent with this id does not exist in runtime")
	}
	
	a.TypeId = agent.TypeId
	tempAgent := &Agent{
		AgentInfo: AgentInfo{
			TypeId:           a.TypeId,
			AgentCoreJsonStr: a.PropJsonStr,
		},
		ac:    ac,
	}

	coreChanged := agent.AgentCoreJsonStr != a.PropJsonStr

	if coreChanged {
		err = tempAgent.loadCore()
		if err != nil {
			return err
		}
	}

	ok, err = models.UpdateAgent(&a)
	if !ok {
		return errors.New("agent with this id does not exist in db")
	}
	if err != nil {
		return err
	}

	agent.Mutex.Lock()
	defer agent.Mutex.Unlock()

	agent.AgentCoreJsonStr = a.PropJsonStr

	if coreChanged {
		if agent.Enable {
			agent.Stop()
		}
		agent.AgentCore = tempAgent.AgentCore
		if a.Enable && a.TypeId == ScheduleAgentId{
			go agent.Run(agent.Ctx, agent, nil, ac.eventHdl.PushEvent)
		}
	} else {
		if agent.Enable != a.Enable {
			if a.Enable && a.TypeId == ScheduleAgentId{
				go agent.Run(agent.Ctx, agent, nil, ac.eventHdl.PushEvent)
			} else {
				agent.Stop()
			}
		}
	}
	agent.Enable = a.Enable
	agent.EventForever = a.EventForever
	agent.EventMaxAge = time.Duration(a.EventMaxAge)

	return nil
}

func (ac *AgentCollection) SetAgentRelation(agentId int, srcs, dsts []int) error {
	agent, ok := ac.agentMap[agentId]
	if !ok {
		return errors.New("no agent with this id")
	}
	err := models.SetAgentRelation(agentId, srcs, dsts)
	if err != nil {
		return err
	}

	agent.Mutex.Lock()
	agent.SrcAgentId = srcs
	agent.DstAgentId = dsts
	agent.Mutex.Unlock()
	return nil
}

func DryRunAgent() {}

func DeleteAgent() {}
