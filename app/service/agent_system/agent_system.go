package agentsystem

type AgentSystem struct {
	Agents   *AgentCollection
	eventHdl *eventHandler
}

func NewAgentSystem(agents *AgentCollection, eventHdl *eventHandler) *AgentSystem {
	agents.eventHdl = eventHdl
	eventHdl.agents = agents
	return &AgentSystem{agents, eventHdl}
}

func (as *AgentSystem) Run() {
	as.eventHdl.run()
	err := as.Agents.init()
	if err != nil {
		panic(err)
	}
}
