package agentsystem

import (
	"fmt"
	"hash/fnv"
	"reflect"
	"sort"

	"github.com/Web-Engineering-XDU/Project-Backend/app/models"
	jsoniter "github.com/json-iterator/go"
)

type eventHandler struct {
	eventChan chan *Event
	agents    *AgentCollection
}

func NewEventHandler() eventHandler {
	return eventHandler{
		eventChan: make(chan *Event, 100),
	}
}

func (eventHdl *eventHandler) run() {
	for i := 0; i < 1; i++ {
		go func() {
			json := jsoniter.ConfigFastest
			var event *Event
			for {
				var err error
				event = <-eventHdl.eventChan

				_, ok := eventHdl.agents.agentMap[event.SrcAgent.ID]
				if !ok {
					continue
				}

				if event.SrcAgent.EventMaxAge != 0 {
					jsonStr := ""
					log := ""
					if event.MetError {
						log = event.Log
						event.ToBeDelivered = false
					} else {
						jsonStr, err = json.MarshalToString(event.Msg)
						if err != nil {
							panic(err)
						}
					}
					event.Msg["$"] = fmt.Sprint(event.SrcAgent.ID)
					eventHash := HashMapString(event.Msg)
					if models.SelectHashCount(eventHash, event.SrcAgent.ID) > 0 && event.SrcAgent.IgnoreDuplicateEvent() {
						event.ToBeDelivered = false
					} else {
						err = models.InsertEvent(&models.Event{
							SrcAgentId:  event.SrcAgent.ID,
							JsonStr:     jsonStr,
							ContentHash: eventHash,
							Error:       event.MetError,
							Log:         log,
							CreateAt:    event.CreateTime,
							DeleteAt:    event.DeleteTime,
						})
						if err != nil {
							panic(err)
						}
					}
				}
				if !event.ToBeDelivered {
					continue
				}
				fmt.Println(event)
				event.SrcAgent.Mutex.RLock()
				for _, v := range event.SrcAgent.DstAgentId {
					go eventHdl.agents.NextAgentDo(v, event)
				}
				event.SrcAgent.Mutex.RUnlock()
			}
		}()
	}
}

func (eventHdl *eventHandler) PushEvent(e *Event) {
	eventHdl.eventChan <- e
}

func HashMapString(m map[string]string) string {
	h := fnv.New64a()
	keys := reflect.ValueOf(m).MapKeys()
	sort.Slice(keys, func(i, j int) bool {
		return keys[i].String() < keys[j].String()
	})
	for _, k := range keys {
		v := reflect.ValueOf(m[k.String()])
		keyBytes := []byte(k.String())
		valueBytes := []byte(v.String())
		h.Write(keyBytes)
		h.Write(valueBytes)
	}
	hash := h.Sum64()
	return fmt.Sprintf("%x", hash)
}
