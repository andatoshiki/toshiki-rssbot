package handler

import (
	"context"
	"strconv"
	"strings"

	"github.com/spf13/cast"
	tb "gopkg.in/telebot.v3"

	"github.com/andatoshiki/toshiki-rssbot/internal/bot/message"
	"github.com/andatoshiki/toshiki-rssbot/internal/bot/session"
	"github.com/andatoshiki/toshiki-rssbot/internal/core"
	"github.com/andatoshiki/toshiki-rssbot/internal/log"
)

type SetUpdateInterval struct {
	core *core.Core
}

func NewSetUpdateInterval(core *core.Core) *SetUpdateInterval {
	return &SetUpdateInterval{core: core}
}

func (s *SetUpdateInterval) Command() string {
	return "/setinterval"
}

func (s *SetUpdateInterval) Description() string {
	return "Configure the update & refresh time interval for subscription feeds"
}

func (s *SetUpdateInterval) getMessageWithoutMention(ctx tb.Context) string {
	mention := message.MentionFromMessage(ctx.Message())
	if mention == "" {
		return ctx.Message().Payload
	}
	return strings.Replace(ctx.Message().Payload, mention, "", -1)
}

func (s *SetUpdateInterval) Handle(ctx tb.Context) error {
	msg := s.getMessageWithoutMention(ctx)
	args := strings.Split(strings.TrimSpace(msg), " ")
	if len(args) < 2 {
		return ctx.Reply("/setinterval [interval] [source_id] Configure the refresh or update interval for subscription feeds, default unit of time is minute (Configuration for multiple sub_id is allowed by splitting with spaces)") 
	}

	interval, err := strconv.Atoi(args[0])
	if interval <= 0 || err != nil {
		return ctx.Reply("Please enter the correct update time interval") 
	}

	subscribeUserID := ctx.Message().Chat.ID
	mentionChat, _ := session.GetMentionChatFromCtxStore(ctx)
	if mentionChat != nil {
		subscribeUserID = mentionChat.ID
	}

	for _, id := range args[1:] {
		sourceID := cast.ToUint(id)
		if err := s.core.SetSubscriptionInterval(
			context.Background(), subscribeUserID, sourceID, interval,
		); err != nil {
			log.Errorf("SetSubscriptionInterval failed, %v", err)
			return ctx.Reply("Failed to configure time intervals!")
		}
	}
	return ctx.Reply("Successfully configured time intervals!")
}

func (s *SetUpdateInterval) Middlewares() []tb.MiddlewareFunc {
	return nil
}
