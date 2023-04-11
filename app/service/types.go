package service

import (
	"context"
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
	Name   string
	Enable bool

	AgentId          int
	AgentType        string
	AgentCoreJsonStr string

	AllowInput, AllowOutput bool
	SrcAgentId, DstAgentId  []int

	EventMaxAge time.Duration
}

type Agent struct {
	AgentInfo
	AgentCore

	EventHdl *eventHandler

	Ctx   context.Context
	Mutex sync.RWMutex
}

type AgentCore interface {
	Run(context.Context, *Agent, *Event)
	Stop()
}
