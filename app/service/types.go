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
	id                      int
	name                    string
	agentType               string
	enable                  bool
	
	allowInput, allowOutput bool
	srcAgentId, dstAgentId  []int

	tempEvent				bool
	eventMaxAge				time.Duration
}

type AgentCore interface {
	Run(ctx context.Context, agent *Agent, event *Event)
	Stop()
}

type Agent struct {
	AgentInfo

	ctx    context.Context
	cancle context.CancelFunc
	mutex  sync.RWMutex

	agentCore *AgentCore
}


