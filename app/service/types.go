package service

import (
	"context"
	"errors"
	"sync"
	"time"
)

type Message map[string]string

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

	EventHdl *eventHandler

	Ctx   context.Context
	Mutex sync.RWMutex
}

func (a *Agent) loadCore() error {
	switch a.AgentInfo.AgentTypeId {
	case 1:
		return a.loadSchduleAgentCore()
	default:
		return errors.New("Unkonw agent core")
	}
	
}

func (a *Agent) loadSchduleAgentCore() error {
	//TODO
	return nil
}


