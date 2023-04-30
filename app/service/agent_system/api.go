package agentsystem

import (
	"sync"

	"github.com/Web-Engineering-XDU/Project-Backend/app/models"
)

func (ac *AgentCollection) AddAgent(a models.Agent) {
	_, ok := ac.agentMap[a.ID]
	if ok {
		//Already running
		return
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
	err := agent.loadCore()
	if err != nil {
		//TODO
		panic(err)
	}
	ac.agentMap[agent.ID] = agent
	if agent.Enable && agent.TypeId == 1 {
		go agent.Run(agent.Ctx, agent, nil)
	}
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

	if agent.AgentCoreJsonStr != a.PropJsonStr {
		agent.Stop()
		agent.AgentCoreJsonStr = a.PropJsonStr
		err := agent.loadCore()
		if err != nil {
			//TODO
			panic(err)
		}
	}
	agent.Enable = a.Enable
	agent.EventForever = a.EventForever
	agent.EventMaxAge = a.EventMaxAge
	agent.Mutex.Unlock()
	
	if agent.Enable && agent.TypeId == 1 {
		go agent.Run(agent.Ctx, agent, nil)
	}
	return true
}

func UpdateRelation() {}

func DryRunAgent() {}

func DeleteAgent() {}
