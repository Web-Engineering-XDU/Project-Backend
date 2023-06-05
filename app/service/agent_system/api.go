package agentsystem

import (
	"context"
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

	err := agent.LoadCore()
	if err != nil {
		return err
	}

	err = models.InsertAgent(a)
	if err != nil {
		return err
	}

	agent.ID = a.ID
	ac.agentMap[agent.ID] = agent
	if agent.TypeId == RssAgentId {
		agent.AgentCore.(*rssAgentCore).loadRssFile(agent)
	}
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
		ac: ac,
	}

	coreChanged := agent.AgentCoreJsonStr != a.PropJsonStr

	if coreChanged {
		err = tempAgent.LoadCore()
		if err != nil {
			return err
		}
	}

	err = models.UpdateAgent(&a)
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
		if a.Enable && a.TypeId == ScheduleAgentId {
			go agent.Run(agent.Ctx, agent, nil, ac.eventHdl.PushEvent)
		}
	} else {
		if agent.Enable != a.Enable {
			if a.Enable && a.TypeId == ScheduleAgentId {
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

	deleted, added := findMissingAndExtraNumbers(agent.SrcAgentId, srcs)

	agent.Mutex.Lock()
	agent.SrcAgentId = srcs
	agent.DstAgentId = dsts
	agent.Mutex.Unlock()

	for _, v := range deleted {
		ac.agentMap[v].Mutex.Lock()
		ac.agentMap[v].DstAgentId = removeFromSlice(ac.agentMap[v].DstAgentId, v)
		ac.agentMap[v].Mutex.Unlock()
	}
	for _, v := range added {
		ac.agentMap[v].Mutex.Lock()
		ac.agentMap[v].DstAgentId = append(ac.agentMap[v].DstAgentId, v)
		ac.agentMap[v].Mutex.Unlock()
	}
	return nil
}

func removeFromSlice(slice []int, num int) []int {
	index := -1

	// 找到要删除的元素的索引
	for i, val := range slice {
		if val == num {
			index = i
			break
		}
	}

	// 如果找到了要删除的元素，则通过切片操作删除它
	if index != -1 {
		slice = append(slice[:index], slice[index+1:]...)
	}

	return slice
}

func findMissingAndExtraNumbers(a, b []int) ([]int, []int) {
	missing := []int{}  // 存储 b 中缺失的数
	extra := []int{}    // 存储 b 中多余的数
	counts := make(map[int]int) // 记录 a 中每个数出现的次数

	// 统计 a 中每个数出现的次数
	for _, num := range a {
		counts[num]++
	}

	// 检查 b 中的每个数在 a 中是否存在
	for _, num := range b {
		if count, found := counts[num]; found && count > 0 {
			counts[num]--  // 在 counts 中减去对应的数的计数
		} else {
			extra = append(extra, num)  // 如果在 counts 中找不到或计数已经为 0，则将其添加到 extra 中
		}
	}

	// 将在 counts 中仍有剩余计数的数添加到 missing 中
	for num, count := range counts {
		if count > 0 {
			missing = append(missing, num)
		}
	}

	return missing, extra
}

func DryRunAgent(a *models.Agent, msg Message) ([]*Message, error) {
	if a.TypeId != HttpAgentId {
		return nil, errors.New("this agent type cannot dry run")
	}
	agent := &Agent{
		AgentInfo: AgentInfo{
			TypeId:           a.TypeId,
			AgentCoreJsonStr: a.PropJsonStr,
		},
		Ctx:   context.Background(),
		Mutex: sync.RWMutex{},
	}
	err := agent.LoadCore()
	if err != nil {
		return nil, err
	}
	var event []*Event
	agent.Run(agent.Ctx, agent, &Event{Msg: msg}, func(e []*Event) {
		event = e
	})
	msgs := make([]*Message, len(event))
	for i, v := range event {
		if v.MetError {
			return nil, errors.New(v.Log)
		}
		msgs[i] = &v.Msg
	}
	return msgs, nil
}

func DeleteAgent() {}
