package agentsystem

import (
	"errors"
	"sync"

	"github.com/Web-Engineering-XDU/Project-Backend/app/models"
)

func (ac *AgentCollection) AddAgent(a models.Agent) error {
	_, ok := ac.agentMap[a.ID]
	if ok {
		//Already running
		return errors.New("agent with this id already exists")
	}
	_, ok = ac.agentTypeMap[a.TypeId]
	if !ok {
		return errors.New("agent type with this id does not exist")
	}

	err := models.InsertAgent(&a)
	if err != nil {
		return err
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
			EventMaxAge:      a.EventMaxAge,
		},
		ac:    ac,
		Ctx:   ac.ctx,
		Mutex: sync.RWMutex{},
	}
	err = agent.loadCore()
	if err != nil {
		//TODO
		panic(err)
	}
	ac.agentMap[agent.ID] = agent
	if agent.Enable && agent.TypeId == 1 {
		go agent.Run(agent.Ctx, agent, nil)
	}
	return nil
}

func (ac *AgentCollection) DeleteAgent(id int) bool {
	agent, ok := ac.agentMap[id]
	if !ok {
		return false
	}

	agent.Stop()
	delete(ac.agentMap, id)
	return models.DeleteAgent(id)
}

func (ac *AgentCollection) UpdateAgent(a models.Agent) bool {
	agent, ok := ac.agentMap[a.ID]
	if !ok {
		return false
	}
	ok = models.UpdateAgent(&a)
	if !ok {
		return false
	}

	agent.Mutex.Lock()

	var err error
	coreChanged := agent.AgentCoreJsonStr != a.PropJsonStr
	agent.AgentCoreJsonStr = a.PropJsonStr

	if coreChanged {
		if agent.Enable {
			agent.Stop()
		}
		err = agent.loadCore()
		if err != nil {
			//TODO
			panic(err)
		}
		if a.Enable {
			go agent.Run(agent.Ctx, agent, nil)
		}
	} else {
		if agent.Enable != a.Enable {
			if a.Enable {
				go agent.Run(agent.Ctx, agent, nil)
			} else {
				agent.Stop()
			}
		}
	}
	agent.Enable = a.Enable
	agent.EventForever = a.EventForever
	agent.EventMaxAge = a.EventMaxAge
	agent.Mutex.Unlock()

	return true
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
