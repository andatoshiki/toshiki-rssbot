package handler

import (
	"fmt"

	tb "gopkg.in/telebot.v3"

	"github.com/andatoshiki/toshiki-rssbot/internal/bot/chat"
	"github.com/andatoshiki/toshiki-rssbot/internal/bot/session"
)

const (
	SetSubscriptionTagButtonUnique = "set_set_sub_tag_btn"
)

type SetSubscriptionTagButton struct {
	bot *tb.Bot
}

func NewSetSubscriptionTagButton(bot *tb.Bot) *SetSubscriptionTagButton {
	return &SetSubscriptionTagButton{bot: bot}
}

func (b *SetSubscriptionTagButton) CallbackUnique() string {
	return "\f" + SetSubscriptionTagButtonUnique
}

func (b *SetSubscriptionTagButton) Description() string {
	return ""
}

func (b *SetSubscriptionTagButton) feedSetAuth(c *tb.Callback, attachData *session.Attachment) bool {
	subscriberID := attachData.GetUserId()
	if subscriberID != c.Sender.ID {
		channelChat, err := b.bot.ChatByID(subscriberID)
		if err != nil {
			return false
		}

		if !chat.IsChatAdmin(b.bot, channelChat, c.Sender.ID) {
			return false
		}
	}
	return true
}

func (b *SetSubscriptionTagButton) Handle(ctx tb.Context) error {
	c := ctx.Callback()
	attachData, err := session.UnmarshalAttachment(ctx.Callback().Data)
	if err != nil {
		return ctx.Edit("System error")
	}

	if !b.feedSetAuth(c, attachData) {
		return ctx.Send("Permission or access rights not granted")
	}
	sourceID := uint(attachData.GetSourceId())
	msg := fmt.Sprintf(
		"Please utilize `/setfeedtag %d tags` command to configure topic tags for the subscription source, `tags` indicates the target tags to be configured. A maximum of 3 tags could be appended to each feed source and tags are required to split by spaces to match the internal grammatical syntax of the bot respectfully \n"+
			"E.g.:`/setfeedtag %d anime moe`",
		sourceID, sourceID,
	)
	return ctx.Edit(msg, &tb.SendOptions{ParseMode: tb.ModeMarkdown})
}

func (b *SetSubscriptionTagButton) Middlewares() []tb.MiddlewareFunc {
	return nil
}
