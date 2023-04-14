package agentsystem

type AgentSystem struct {
	agents   *agentCollection
	eventHdl *eventHandler
}

func NewAgentSystem(agents *agentCollection, eventHdl *eventHandler) *AgentSystem {
	agents.eventHdl = eventHdl
	eventHdl.agents = agents
	return &AgentSystem{agents, eventHdl}
}

func (as *AgentSystem) Run() {
	as.eventHdl.run()
	err := as.agents.init()
	if err != nil {
		panic(err)
	}
	select{}
}
