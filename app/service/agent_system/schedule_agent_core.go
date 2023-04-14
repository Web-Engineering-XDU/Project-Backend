package agentsystem

import (
	"context"
	"log"
	"time"

	"github.com/robfig/cron/v3"
)

var cronTimer = cron.New(cron.WithSeconds())

type ScheduleAgentCore struct {
	CronSpec    string	`json:"cron"`
	cronEntryID cron.EntryID
}

func (sac *ScheduleAgentCore) Run(ctx context.Context, agent *Agent, event *Event) {
	go cronTimer.Run()
	var err error
	sac.cronEntryID, err = cronTimer.AddFunc(sac.CronSpec, func() {
		agent.ac.eventHdl.PushEvent(&Event{
			SrcAgent: agent,
			CreateTime: time.Now(),
			DeleteTime: time.Now().Add(agent.EventMaxAge),
			Msg:      emptyMsg,
		})
	})
	if err != nil {
		log.Panicln(err)
		//TODO 
	}
}

func (sac *ScheduleAgentCore) Stop() {
	if cronTimer != nil{
		cronTimer.Remove(sac.cronEntryID)
	}
}
