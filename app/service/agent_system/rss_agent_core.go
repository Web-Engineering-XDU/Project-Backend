package agentsystem

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/gorilla/feeds"
	"github.com/mmcdole/gofeed"
)

type rssItemTemplate struct {
	Title       string
	Link        string
	Description string
	Author      string
}

func (item rssItemTemplate) render(bindings map[string]string) *rssItemTemplate {
	var err error
	res := make(map[string]interface{}, len(bindings))
	for k, v := range bindings {
		res[k] = v
	}
	item.Title, err = engine.ParseAndRenderString(item.Title, res)
	if err != nil {
		panic(err)
	}
	item.Link, err = engine.ParseAndRenderString(item.Link, res)
	if err != nil {
		panic(err)
	}
	item.Description, err = engine.ParseAndRenderString(item.Description, res)
	if err != nil {
		panic(err)
	}
	item.Author, err = engine.ParseAndRenderString(item.Author, res)
	if err != nil {
		panic(err)
	}
	return &item
}

type rssAgentCore struct {
	Title       string
	Link        string
	Description string
	Author      string
	Template    rssItemTemplate

	file *os.File
	feed *feeds.Feed
}

func (rac *rssAgentCore) loadRssFile(a *Agent) {
	var err error
	a.Mutex.Lock()
	if rac.file == nil {
		err = os.Mkdir("./rss", 0777)
		if err != nil && !os.IsExist(err) {
			panic(err)
		}
		rac.file, err = os.OpenFile(fmt.Sprintf("./rss/%v.xml", a.ID), os.O_RDWR|os.O_CREATE, 0777)
		if err != nil {
			panic(err)
		}
	}
	if rac.feed == nil {
		fileInfo, err := rac.file.Stat()
		if err != nil {
			panic(err)
		}
		content := make([]byte, fileInfo.Size())
		_, err = rac.file.Read(content)
		if err != nil {
			panic(err)
		}
		rac.feed = &feeds.Feed{
			Title:       rac.Title,
			Link:        &feeds.Link{Href: rac.Link},
			Description: rac.Description,
			Author:      &feeds.Author{Name: rac.Author},
		}
		buffer := bytes.NewBuffer(content)
		feed, err := gofeed.NewParser().Parse(buffer)
		if err != nil {
			rac.feed.Created = time.Now()
			rac.feed.Items = []*feeds.Item{}
		} else {
			rac.feed.Created = (*feed.PublishedParsed).Local()
			rac.feed.Items = make([]*feeds.Item, 0, len(feed.Items))
			for _, v := range feed.Items {
				rac.feed.Items = append(rac.feed.Items, &feeds.Item{
					Title:       v.Title,
					Link:        &feeds.Link{Href: v.Link},
					Description: v.Description,
					Author:      &feeds.Author{Name: v.Author.Name},
					Created:     (*v.PublishedParsed).Local(),
				})
			}
		}
	}
	if !a.Enable {
		a.Stop()
	}
	a.Mutex.Unlock()
}

func (a *Agent) loadRssAgentCore() error {
	core := &rssAgentCore{}
	err := json.UnmarshalFromString(a.AgentCoreJsonStr, core)
	if err != nil {
		return err
	}
	a.AgentCore = core
	core.loadRssFile(a)
	return nil
}

func (rac *rssAgentCore) Run(ctx context.Context, agent *Agent, event *Event, callBack func(e []*Event)) {
	var err error
	rac.loadRssFile(agent)
	newItem := rac.Template.render(event.Msg)

	agent.Mutex.Lock()
	rac.feed.Items = append(rac.feed.Items, &feeds.Item{
		Title:       newItem.Title,
		Link:        &feeds.Link{Href: newItem.Link},
		Description: newItem.Description,
		Author:      &feeds.Author{Name: newItem.Author},
		Created:     time.Now(),
	})
	sort.Slice(rac.feed.Items, func(i, j int) bool {
		return rac.feed.Items[i].Created.After(rac.feed.Items[j].Created)
	})
	agent.Mutex.Unlock()

	rss, err := rac.feed.ToRss()
	if err != nil {
		panic(err)
	}
	err = rac.file.Truncate(int64(len(rss)))
	if err != nil {
		panic(err)
	}
	_, err = rac.file.WriteAt([]byte(rss), 0)
	if err != nil {
		panic(err)
	}
}

func (rac *rssAgentCore) Stop() {
	if rac.file != nil {
		rac.file.Close()
		rac.file = nil
	}
	rac.feed = nil
}

func (*rssAgentCore) IgnoreDuplicateEvent() bool {
	return true
}

func (*rssAgentCore) ValidCheck() error {
	return nil
}
