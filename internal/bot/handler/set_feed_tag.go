package handler

import (
	"context"
	"strings"

	"github.com/spf13/cast"
	tb "gopkg.in/telebot.v3"

	"github.com/andatoshiki/toshiki-rssbot/internal/bot/message"
	"github.com/andatoshiki/toshiki-rssbot/internal/bot/session"
	"github.com/andatoshiki/toshiki-rssbot/internal/core"
)

type SetFeedTag struct {
	core *core.Core
}

func NewSetFeedTag(core *core.Core) *SetFeedTag {
	return &SetFeedTag{core: core}
}

func (s *SetFeedTag) Command() string {
	return "/setfeedtag"
}

func (s *SetFeedTag) Description() string {
	return "Configure and set tags for RSS subscriptions"
}

func (s *SetFeedTag) getMessageWithoutMention(ctx tb.Context) string {
	mention := message.MentionFromMessage(ctx.Message())
	if mention == "" {
		return ctx.Message().Payload
	}
	return strings.Replace(ctx.Message().Payload, mention, "", -1)
}

func (s *SetFeedTag) Handle(ctx tb.Context) error {
	msg := s.getMessageWithoutMention(ctx)
	args := strings.Split(strings.TrimSpace(msg), " ")
	if len(args) < 1 {
		return ctx.Reply("/setfeedtag `[source_id] [tag1] [tag2]` to set tags for subscription feeds; a maximum of 3 tags could be appended to each feed source and tags are required to split by spaces to match the internal grammatical syntax of the bot respectfully")
	}

	// truncate properties
	if len(args) > 4 {
		args = args[:4]
	}

	sourceID := cast.ToUint(args[0])
	mentionChat, _ := session.GetMentionChatFromCtxStore(ctx)
	subscribeUserID := ctx.Chat().ID
	if mentionChat != nil {
		subscribeUserID = mentionChat.ID
	}

	if err := s.core.SetSubscriptionTag(context.Background(), subscribeUserID, sourceID, args[1:]); err != nil {
		return ctx.Reply("Failed to set tag(s) for subscription feed!")
	}
	return ctx.Reply("Successfully set tag(s) for subscription feed!")
}

func (s *SetFeedTag) Middlewares() []tb.MiddlewareFunc {
	return nil
}
