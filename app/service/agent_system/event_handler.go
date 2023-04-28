package agentsystem

import (
	"fmt"
	"hash/fnv"
	"reflect"
	"sort"

	"github.com/Web-Engineering-XDU/Project-Backend/app/models"
	"github.com/go-sql-driver/mysql"
	jsoniter "github.com/json-iterator/go"
)

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
	for i := 0; i < 1; i++ {
		go func() {
			json := jsoniter.ConfigFastest
			var event *Event
			for {
				var err error
				event = <-eventHdl.eventChan
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
					err = models.InsertEvent(&models.Event{
						SrcAgentId:  event.SrcAgent.Id,
						JsonStr:     jsonStr,
						ContentHash: HashMapString(event.Msg),
						Error:       event.MetError,
						Log:         log,
						CreateAt:    event.CreateTime,
						DeleteAt:    event.DeleteTime,
					})
					if err != nil {
						switch err.(*mysql.MySQLError).Number {
						case 1062:
							event.ToBeDelivered = false
						default:
							panic(err)
						}
					}
				}
				if !event.ToBeDelivered {
					continue
				}
				fmt.Println(event)
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
