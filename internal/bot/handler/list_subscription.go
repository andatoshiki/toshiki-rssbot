package handler

import (
	"context"
	"fmt"
	"strings"

	tb "gopkg.in/telebot.v3"

	"github.com/andatoshiki/toshiki-rssbot/internal/bot/chat"
	"github.com/andatoshiki/toshiki-rssbot/internal/bot/message"
	"github.com/andatoshiki/toshiki-rssbot/internal/core"
	"github.com/andatoshiki/toshiki-rssbot/internal/log"
	"github.com/andatoshiki/toshiki-rssbot/internal/model"
)

const (
	MaxSubsSizePerPage = 50
)

type ListSubscription struct {
	core *core.Core
}

func NewListSubscription(core *core.Core) *ListSubscription {
	return &ListSubscription{core: core}
}

func (l *ListSubscription) Command() string {
	return "/list"
}

func (l *ListSubscription) Description() string {
	return "All subscribed rss source feeds"
}

func (l *ListSubscription) listChatSubscription(ctx tb.Context) error {
	// private chat or group
	if ctx.Chat().Type != tb.ChatPrivate && !chat.IsChatAdmin(ctx.Bot(), ctx.Chat(), ctx.Sender().ID) {
		return ctx.Send("No permission or access rights")
	}

	stdCtx := context.Background()
	sources, err := l.core.GetUserSubscribedSources(stdCtx, ctx.Chat().ID)
	if err != nil {
		log.Errorf("GetUserSubscribedSources failed, %v", err)
		return ctx.Send("Errors or failures to fetch subscription list")
	}

	return l.replaySubscribedSources(ctx, sources)
}

func (l *ListSubscription) listChannelSubscription(ctx tb.Context, channelName string) error {
	channelChat, err := ctx.Bot().ChatByUsername(channelName)
	if err != nil {
		return ctx.Send("Failed to fetch channel information")
	}

	if !chat.IsChatAdmin(ctx.Bot(), channelChat, ctx.Sender().ID) {
		return ctx.Send("Bot operational executions by non-administrative users of chanel are not permitted") 
	}

	stdCtx := context.Background()
	sources, err := l.core.GetUserSubscribedSources(stdCtx, channelChat.ID)
	if err != nil {
		log.Errorf("GetUserSubscribedSources failed, %v", err)
		return ctx.Send("Errors or failures to fetch subscriptions")
	}
	return l.replaySubscribedSources(ctx, sources)
}

func (l *ListSubscription) Handle(ctx tb.Context) error {
	mention := message.MentionFromMessage(ctx.Message())
	if mention != "" {
		return l.listChannelSubscription(ctx, mention)
	}
	return l.listChatSubscription(ctx)
}

func (l *ListSubscription) Middlewares() []tb.MiddlewareFunc {
	return nil
}

func (l *ListSubscription) replaySubscribedSources(ctx tb.Context, sources []*model.Source) error {
	if len(sources) == 0 {
		return ctx.Send("The current subscription list is empty") 
	}
	var msg strings.Builder
	msg.WriteString(fmt.Sprintf("Total of %d subscription(s), subscription list\n", len(sources)))
	count := 0
	for i := range sources {
		msg.WriteString(fmt.Sprintf("[[%d]] [%s](%s)\n", sources[i].ID, sources[i].Title, sources[i].Link))
		count++
		if count == MaxSubsSizePerPage {
			ctx.Send(msg.String(), &tb.SendOptions{DisableWebPagePreview: true, ParseMode: tb.ModeMarkdown})
			count = 0
			msg.Reset()
		}
	}

	if count != 0 {
		ctx.Send(msg.String(), &tb.SendOptions{DisableWebPagePreview: true, ParseMode: tb.ModeMarkdown})
	}
	return nil
}
