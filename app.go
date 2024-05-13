package main

import (
	"MrRSS/backend"

	"context"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) InitDatabase() {
	backend.InitDatabase()
	feeds := []backend.FeedsInfo{
		{Link: "https://www.kawabangga.com/feed", Category: "blog"},
		{Link: "https://jvns.ca/atom.xml", Category: "wensite"},
		{Link: "https://www.ruanyifeng.com/blog/atom.xml", Category: "blog"},
		{Link: "https://www.appinn.com/feed/", Category: "website"},
	}
	backend.SetFeedList(feeds)
}

func (a *App) GetFeedList() []backend.FeedsInfo {
	return backend.GetFeedList()
}

func (a *App) GetFeedContent() []backend.FeedContentsInfo {
	return backend.GetFeedContent()
}

func (a *App) GetHistory() []backend.FeedContentsInfo {
	return backend.GetHistory()
}

func (a *App) SetHistory(history []backend.FeedContentsInfo) {
	backend.SetHistory(history)
}

func (a *App) SetHistoryReaded(feed backend.FeedContentsInfo) {
	backend.SetHistoryReaded(feed)
}

func (a *App) ClearHistory() {
	backend.ClearHistory()
}
