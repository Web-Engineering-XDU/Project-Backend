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

	MetError      bool
	Log           string
	ToBeDelivered bool
}

type AgentTypeProp struct {
	AllowInput, AllowOutput bool
}

type AgentInfo struct {
	ID     int
	Enable bool

	TypeId           int
	AgentCoreJsonStr string

	SrcAgentId, DstAgentId []int

	EventForever bool
	EventMaxAge  time.Duration
}

type AgentCore interface {
	Run(context.Context, *Agent, *Event)
	Stop()
}

type Agent struct {
	AgentInfo
	AgentCore

	ac *AgentCollection

	Ctx   context.Context
	Mutex sync.RWMutex
}

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func (a *Agent) loadCore() error {
	switch a.TypeId {
	case 1:
		return a.loadSchduleAgentCore()
	case 2:
		return a.loadPrintAgentCore()
	case 3:
		return a.loadHttpAgentCore()
	default:
		return errors.New("unkonw agent core")
	}

}
