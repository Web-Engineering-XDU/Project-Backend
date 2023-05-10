package agentsystem

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	jsoniter "github.com/json-iterator/go"
)

const (
	ScheduleAgentId = iota + 1 //1
	HttpAgentId                //2
	PrintAgentId               //3
	RssAgentId                 //4
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
	Run(context.Context, *Agent, *Event, func(e *Event))
	Stop()
	IgnoreDuplicateEvent() bool
	ValidCheck() error
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
	var err error
	switch a.TypeId {
	case ScheduleAgentId:
		err = a.loadSchduleAgentCore()
	case PrintAgentId:
		err = a.loadPrintAgentCore()
	case HttpAgentId:
		err = a.loadHttpAgentCore()
	case RssAgentId:
		err = a.loadHttpAgentCore()
	default:
		err = errors.New("unkonw agent core")
	}
	if err != nil {
		return errors.New(fmt.Sprintf("load core(id = %v) error: %v", a.TypeId, err.Error()))
	}

	return a.ValidCheck()
}
