package service

import (
	"context"
	"sync"
	"time"
)

type Message map[string]string

type Event struct {
	SrcAgent *Agent
	Msg      Message
}

type AgentInfo struct {
	Id        int
	Name      string
	AgentType string
	Enable    bool

	AllowInput, AllowOutput bool
	SrcAgentId, DstAgentId  []int

	TempEvent   bool
	EventMaxAge time.Duration
}

type Agent struct {
	AgentInfo

	EventHdl *eventHandler

	Ctx   context.Context
	Mutex sync.RWMutex

	Run  func(context.Context, *Event)
	Stop func()
}
