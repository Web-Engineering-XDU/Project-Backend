package agentsystem

import (
	"context"
	"errors"
	"sync"
	"time"

	jsoniter "github.com/json-iterator/go"
)

type Message map[string]interface{}

var emptyMsg = map[string]interface{}{}

type Event struct {
	SrcAgent   *Agent
	CreateTime time.Time
	DeleteTime time.Time
	Msg        Message

	MetError      bool
	Log           string
	ToBeDelivered bool
}

type AgentInfo struct {
	Id     int
	Enable bool

	AgentTypeId      int
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
