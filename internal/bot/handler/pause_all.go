package handler

import (
	"context"
	"fmt"

	tb "gopkg.in/telebot.v3"

	"github.com/andatoshiki/toshiki-rssbot/internal/bot/session"
	"github.com/andatoshiki/toshiki-rssbot/internal/core"
)

type PauseAll struct {
	core *core.Core
}

func NewPauseAll(core *core.Core) *PauseAll {
	return &PauseAll{core: core}
}

func (p *PauseAll) Command() string {
	return "/pauseall"
}

func (p *PauseAll) Description() string {
	return "Pause and terminate fetching all subscription updates"
}

func (p *PauseAll) Handle(ctx tb.Context) error {
	subscribeUserID := ctx.Message().Chat.ID
	var channelChat *tb.Chat
	v := ctx.Get(session.StoreKeyMentionChat.String())
	if v != nil {
		var ok bool
		channelChat, ok = v.(*tb.Chat)
		if ok && channelChat != nil {
			subscribeUserID = channelChat.ID
		}
	}

	source, err := p.core.GetUserSubscribedSources(context.Background(), subscribeUserID)
	if err != nil {
		return ctx.Reply("Internal system error")
	}

	for _, s := range source {
		err := p.core.DisableSourceUpdate(context.Background(), s.ID)
		if err != nil {
			return ctx.Reply("Failed to pause")
		}
	}

	reply := "All subscription updates have been paused and terminated"
	if channelChat != nil {
		reply = fmt.Sprintf("All subscriptions of channel [%s](https://t.me/%s) have been completely paused and terminated", channelChat.Title, channelChat.Username)
	}
	return ctx.Send(
		reply, &tb.SendOptions{
			DisableWebPagePreview: true,
			ParseMode:             tb.ModeMarkdown,
		},
	)
}

func (p *PauseAll) Middlewares() []tb.MiddlewareFunc {
	return nil
}
