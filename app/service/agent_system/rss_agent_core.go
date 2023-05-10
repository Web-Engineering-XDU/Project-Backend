package agentsystem

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/Web-Engineering-XDU/Project-Backend/config"
	. "github.com/gorilla/feeds"
	"github.com/ilyakaznacheev/cleanenv"
)

type RssAgentCore struct {
}

func (a *Agent) loadRssAgentCore() error {
	core := &RssAgentCore{}
	a.AgentCore = core
	return nil
}

func (pac *RssAgentCore) Run(ctx context.Context, agent *Agent, event *Event, callBack func(e *Event)) {
	now := time.Now()
	feed := &Feed{
		Title:       "Huggo event ",
		Description: "discussion about tech, footie, photos",
		Author:      &Author{Name: "Jason Moiron", Email: "jmoiron@jmoiron.net"},
		Updated:     now,
	}
	feed.Items = []*Item{
		{
			Title:       `${agent.ID} Recive Event from ${event.SrcAgent.ID}", , event.Msg, `,
			Description: `${event.Msg}`,
			Created:     now,
		},
	}

	// TODO: Need better way to read the config.
	ex, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}

	var config config.Config
	err = cleanenv.ReadConfig(filepath.Dir(ex)+"/config.yml", &config)
	if err != nil {
		log.Fatal(err)
	}

	// TODO:
	newFile, err := os.Create(config.RssPath.Path + "index.xml")
	if err != nil {
		log.Fatal(err)
	}

	feed.WriteRss(newFile)

	newFile.Close()

	//fmt.Printf("%v Recive Event: %v from %v\n", agent.ID, event.Msg, event.SrcAgent.ID)
}

func (*RssAgentCore) Stop() {}

func (*RssAgentCore) IgnoreDuplicateEvent() bool {
	return true
}

func (*RssAgentCore) ValidCheck() error {
	return nil
}
