package agentsystem

import (
	"context"
	"fmt"
)

type PrintAgentCore struct{}

func (pac *PrintAgentCore) Run(ctx context.Context, agent *Agent, event *Event) {
	fmt.Printf("%v Recive Event: %v from %v\n", agent.Id ,event.Msg, event.SrcAgent.Id)
}

func (*PrintAgentCore) Stop() {}