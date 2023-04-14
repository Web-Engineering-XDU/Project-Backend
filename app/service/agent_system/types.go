package agentsystem

import (
	"context"
	"errors"
	"sync"
	"time"

	jsoniter "github.com/json-iterator/go"
)

type Message map[string]string
var emptyMsg = map[string]string{}

type Event struct {
	SrcAgent   *Agent
	CreateTime time.Time
	DeleteTime time.Time
	Msg        Message
}

type AgentInfo struct {
	Id     int
	Enable bool

	AgentTypeId          int
	AgentCoreJsonStr string

	AllowInput, AllowOutput bool
	SrcAgentId, DstAgentId  []int

	EventMaxAge time.Duration
}

type AgentCore interface {
	Run(context.Context, *Agent, *Event)
	Stop()
}

type Agent struct {
	AgentInfo
	AgentCore

	ac *agentCollection

	Ctx   context.Context
	Mutex sync.RWMutex
}

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func (a *Agent) loadCore() error {
	switch a.AgentInfo.AgentTypeId {
	case 1:
		return a.loadSchduleAgentCore()
	case 2:
		return a.loadPrintAgentCore()
	default:
		return errors.New("unkonw agent core")
	}
	
}

func (a *Agent) loadSchduleAgentCore() error {
	core := &ScheduleAgentCore{}
	err := json.UnmarshalFromString(a.AgentCoreJsonStr, core)
	if err != nil {
		return err
	}
	a.AgentCore = core
	return nil
}

func (a *Agent) loadPrintAgentCore() error {
	core := &PrintAgentCore{}
	a.AgentCore = core
	return nil
}

