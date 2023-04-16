package agentsystem

type eventHandler struct {
	eventChan chan *Event
	agents    *agentCollection
}

func NewEventHandler() eventHandler {
	return eventHandler{
		eventChan: make(chan *Event, 100),
	}
}

func (eventHdl *eventHandler) run() {
	for i := 0; i < 10; i++ {
		go func() {
			var event *Event
			for {
				event = <-eventHdl.eventChan
				if event.MetError {
					//TODO
				}
				if !event.ToBeDelivered {
					return
				}
				for _, v := range event.SrcAgent.DstAgentId {
					go eventHdl.agents.NextAgentDo(v, event)
				}
			}
		}()
	}
}

func (eventHdl *eventHandler) PushEvent(e *Event) {
	eventHdl.eventChan <- e
}
