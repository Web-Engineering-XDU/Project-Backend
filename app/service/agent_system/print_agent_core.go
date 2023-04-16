package agentsystem

import (
	"context"
	"fmt"
)

type printAgentCore struct{}

func (a *Agent) loadPrintAgentCore() error {
	core := &printAgentCore{}
	a.AgentCore = core
	return nil
}

func (pac *printAgentCore) Run(ctx context.Context, agent *Agent, event *Event) {
	fmt.Printf("%v Recive Event: %v from %v\n", agent.Id, event.Msg, event.SrcAgent.Id)
}

func (*printAgentCore) Stop() {}
