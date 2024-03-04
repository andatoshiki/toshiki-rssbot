package handler

import (
	"bytes"
	"context"
	"text/template"

	tb "gopkg.in/telebot.v3"

	"github.com/andatoshiki/toshiki-rssbot/internal/bot/chat"
	"github.com/andatoshiki/toshiki-rssbot/internal/bot/session"
	"github.com/andatoshiki/toshiki-rssbot/internal/config"
	"github.com/andatoshiki/toshiki-rssbot/internal/core"
)

const (
	SubscriptionSwitchButtonUnique = "set_toggle_update_btn"
)

type SubscriptionSwitchButton struct {
	bot  *tb.Bot
	core *core.Core
}

func NewSubscriptionSwitchButton(bot *tb.Bot, core *core.Core) *SubscriptionSwitchButton {
	return &SubscriptionSwitchButton{bot: bot, core: core}
}

func (b *SubscriptionSwitchButton) CallbackUnique() string {
	return "\f" + SubscriptionSwitchButtonUnique
}

func (b *SubscriptionSwitchButton) Description() string {
	return ""
}

func (b *SubscriptionSwitchButton) Handle(ctx tb.Context) error {
	c := ctx.Callback()
	if c == nil {
		return ctx.Respond(&tb.CallbackResponse{Text: "error"})
	}

	attachData, err := session.UnmarshalAttachment(ctx.Callback().Data)
	subscriberID := attachData.GetUserId()
	if subscriberID != c.Sender.ID {
		// If the subscriber ID is different from the button user's ID, administrator permission needs to be verified.
		channelChat, err := b.bot.ChatByID(subscriberID)
		if err != nil {
			return ctx.Respond(&tb.CallbackResponse{Text: "error"})
		}
		if !chat.IsChatAdmin(b.bot, channelChat, c.Sender.ID) {
			return ctx.Respond(&tb.CallbackResponse{Text: "error"})
		}
	}

	sourceID := uint(attachData.GetSourceId())
	sub, err := b.core.GetSubscription(context.Background(), subscriberID, sourceID)
	if sub == nil || err != nil {
		return ctx.Respond(&tb.CallbackResponse{Text: "error"})
	}

	err = b.core.ToggleSourceUpdateStatus(context.Background(), sourceID)
	if err != nil {
		return ctx.Respond(&tb.CallbackResponse{Text: "error"})
	}

	source, _ := b.core.GetSource(context.Background(), sourceID)
	t := template.New("setting template")
	_, _ = t.Parse(feedSettingTmpl)

	text := new(bytes.Buffer)
	_ = t.Execute(text, map[string]interface{}{"source": source, "sub": sub, "Count": config.ErrorThreshold})
	_ = ctx.Respond(&tb.CallbackResponse{Text: "Successfully modified"})
	return ctx.Edit(
		text.String(),
		&tb.SendOptions{ParseMode: tb.ModeHTML},
		&tb.ReplyMarkup{InlineKeyboard: genFeedSetBtn(c, sub, source)},
	)
}

func (b *SubscriptionSwitchButton) Middlewares() []tb.MiddlewareFunc {
	return nil
}
