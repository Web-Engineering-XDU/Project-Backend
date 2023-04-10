package service

type eventHandler struct {
	eventChan chan *Event
	agents    *agentCollection
}

func newEventHandler() eventHandler {
	return eventHandler{
		eventChan: make(chan *Event, 100),
	}
}

func (eventHdl *eventHandler) startEventHandler() {
	for i := 0; i < 10; i++ {
		go func() {
			var event *Event
			for {
				event = <-eventHdl.eventChan
				event.SrcAgent.Mutex.RLock()
				for _, v := range event.SrcAgent.DstAgentId {
					go eventHdl.agents.NextAgentDo(v, event)
				}
				event.SrcAgent.Mutex.Unlock()
			}
		}()
	}
}

func (eventHdl *eventHandler) PushEvent(e *Event) {
	eventHdl.eventChan <- e
}
