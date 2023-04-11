package agents

import (
	"context"
	"time"

	"github.com/Web-Engineering-XDU/Project-Backend/app/service"
	"github.com/robfig/cron/v3"
)

var cronTimer = cron.New()

type ScheduleAgentCore struct {
	cronSpec    string
	cronEntryID cron.EntryID
}

func (sac *ScheduleAgentCore) Run(ctx context.Context, agent *service.Agent, event *service.Event) {
	cronTimer.Run()
	var err error
	sac.cronEntryID, err = cronTimer.AddFunc(sac.cronSpec, func() {
		agent.EventHdl.PushEvent(&service.Event{
			SrcAgent: agent,
			CreateTime: time.Now(),
			DeleteTime: time.Now().Add(agent.EventMaxAge),
			Msg:      nil,
		})
	})
	if err != nil {
		//TODO
	}
}

func (sac *ScheduleAgentCore) Stop() {
	if cronTimer != nil {
		cronTimer.Remove(sac.cronEntryID)
	}
}
