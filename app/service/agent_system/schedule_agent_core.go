package agentsystem

import (
	"context"
	"log"
	"time"

	"github.com/robfig/cron/v3"
)

var cronTimer = cron.New(cron.WithSeconds())

type scheduleAgentCore struct {
	CronSpec    string `json:"cron"`
	cronEntryID cron.EntryID
}

func (a *Agent) loadSchduleAgentCore() error {
	core := &scheduleAgentCore{}
	err := json.UnmarshalFromString(a.AgentCoreJsonStr, core)
	if err != nil {
		return err
	}
	a.AgentCore = core
	return nil
}

func (sac *scheduleAgentCore) Run(ctx context.Context, agent *Agent, event *Event, callBack func(e *Event)) {
	go cronTimer.Run()
	var err error
	agent.Mutex.RLock()
	sac.cronEntryID, err = cronTimer.AddFunc(sac.CronSpec, func() {
		callBack(&Event{
			SrcAgent:   agent,
			CreateTime: time.Now(),
			DeleteTime: time.Now().Add(agent.EventMaxAge),
			Msg:        emptyMsg,
			ToBeDelivered: true,
		})
	})
	agent.Mutex.RUnlock()
	if err != nil {
		log.Panicln(err)
	}
}

func (sac *scheduleAgentCore) Stop() {
	if cronTimer != nil {
		cronTimer.Remove(sac.cronEntryID)
	}
}

func (sac *scheduleAgentCore) IgnoreDuplicateEvent() bool {
	return true
}

func (sac *scheduleAgentCore) ValidCheck() error {
	p := cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	_, err := p.Parse(sac.CronSpec)
	return err
}

