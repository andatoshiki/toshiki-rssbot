package handler

import (
	"context"
	"fmt"

	tb "gopkg.in/telebot.v3"

	"github.com/andatoshiki/toshiki-rssbot/internal/bot/session"
	"github.com/andatoshiki/toshiki-rssbot/internal/core"
)

type ActiveAll struct {
	core *core.Core
}

func NewActiveAll(core *core.Core) *ActiveAll {
	return &ActiveAll{core: core}
}

func (a *ActiveAll) Command() string {
	return "/activeall"
}

func (a *ActiveAll) Description() string {
	return "Start fetching subscription updates"
}

func (a *ActiveAll) Handle(ctx tb.Context) error {
	mentionChat, _ := session.GetMentionChatFromCtxStore(ctx)
	subscribeUserID := ctx.Chat().ID
	if mentionChat != nil {
		subscribeUserID = mentionChat.ID
	}

	source, err := a.core.GetUserSubscribedSources(context.Background(), subscribeUserID)
	if err != nil {
		return ctx.Reply("Internal service error")
	}

	for _, s := range source {
		err := a.core.EnableSourceUpdate(context.Background(), s.ID)
		if err != nil {
			return ctx.Reply("Activation failed")
		}
	}

	reply := "All subscriptions has been enabled and activated"
	if mentionChat != nil {
		reply = fmt.Sprintf("Channel [%s](https://t.me/%s) has enabled and activated all feed subscription sources update", mentionChat.Title, mentionChat.Username)
	}

	return ctx.Reply(
		reply, &tb.SendOptions{
			DisableWebPagePreview: true,
			ParseMode:             tb.ModeMarkdown,
		},
	)
}

func (a *ActiveAll) Middlewares() []tb.MiddlewareFunc {
	return nil
}
