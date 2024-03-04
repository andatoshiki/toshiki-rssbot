package handler

import (
	tb "gopkg.in/telebot.v3"
)

type Help struct {
}

func NewHelp() *Help {
	return &Help{}
}

func (h *Help) Command() string {
	return "/help"
}

func (h *Help) Description() string {
	return "Help"
}

func (h *Help) Handle(ctx tb.Context) error {
	message := `
	/sub Subscribe an RSS feed source to your feed list
	/unsub  Remove a subscription source from your existing feed list
	/list View all existing subscription sources
	/set Configure & manage subscription list
	/check Inspect the existing subscribed feed list status
	/setfeedtag Append a custom tag to a subscription source
	/setinterval Configure the refresh interval for a subscription source
	/activeall Resume & enable all existing subscription sources
	/pauseall Pause & terminate all existing subscription sources
	/help View help & support information
	/import Import your subscription list to an OPML file
	/export Export your subscription list to an OPML file
	/unsuball Remove and cancel all existing subscriptions
	Visit for more detailed bot usage & affiliated documentation at https://note.toshiki.dev/
	Made with love and coffee by @andatoshiki at @toshikidev in Arizona State University & proudly open sourced on GitLab
	`
	return ctx.Send(message)
}

func (h *Help) Middlewares() []tb.MiddlewareFunc {
	return nil
}
