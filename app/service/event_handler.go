package service

type eventHandler struct {
	eventChan chan *Event
}

var eh = eventHandler{
	eventChan: make(chan *Event, 100),
}

func GetEventHandler() eventHandler {
	return eh
}

func StartEventHandler() {
	for i := 0; i < 10; i++ {
		go func() {
			var event *Event
			for {
				event = <-eh.eventChan
				event.SrcAgent.mutex.RLock()
				for _, v := range event.SrcAgent.dstAgentId {
					go NextAgentDo(v, event)
				}
				event.SrcAgent.mutex.Unlock()
			}
		}()
	}
}

func PushEvent(e *Event) {
	eh.eventChan <- e
}