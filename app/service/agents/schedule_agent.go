package agents

import (
	"context"

	"github.com/Web-Engineering-XDU/Project-Backend/app/service"
	"github.com/robfig/cron/v3"
)

var cronTimer = cron.New()

type ScheduleAgent struct {
	service.Agent

	cronSpec    string
	cronEntryID cron.EntryID
}

func (sac *ScheduleAgent) Run(ctx context.Context, event *service.Event) {
	cronTimer.Run()
	var err error
	sac.cronEntryID, err = cronTimer.AddFunc(sac.cronSpec, func() {
		sac.EventHdl.PushEvent(&service.Event{
			SrcAgent: &sac.Agent,
			Msg:      nil,
		})
	})
	if err != nil {
		//TODO
	}
}

func (sac *ScheduleAgent) Stop() {
	if cronTimer != nil {
		cronTimer.Remove(sac.cronEntryID)
	}
}
