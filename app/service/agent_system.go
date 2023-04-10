package service

type AgentSystem struct {
	agents   *agentCollection
	eventHdl *eventHandler
}

func newAgentSystem(agents *agentCollection, eventHdl *eventHandler) *AgentSystem {
	agents.eventHdl = eventHdl
	eventHdl.agents = agents
	return &AgentSystem{agents, eventHdl}
}

// func StartService() {
// 	startAgentManager()
// }