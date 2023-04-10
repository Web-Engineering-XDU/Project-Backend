package agents

import (
	"context"

	"github.com/Web-Engineering-XDU/Project-Backend/app/service"
	"github.com/robfig/cron/v3"
)

type ScheduleAgentCore struct {
	cronTimer *cron.Cron
	cronSpec  string

	cronEntryID cron.EntryID
}

func (sac *ScheduleAgentCore) Run(ctx context.Context, agent *service.Agent, event *service.Event) {
	if sac.cronTimer == nil {
		sac.cronTimer = cron.New()
		sac.cronTimer.Run()
	}
	var err error
	sac.cronEntryID, err = sac.cronTimer.AddFunc(sac.cronSpec, func() {
		service.PushEvent(&service.Event{
			SrcAgent: agent,
			Msg:      nil,
		})
	})
	if err != nil {
		//TODO
	}
}

func (sac *ScheduleAgentCore) Stop() {
	if sac.cronTimer != nil {
		sac.cronTimer.Remove(sac.cronEntryID)
	}
}
