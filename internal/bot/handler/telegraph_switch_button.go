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
	TelegraphSwitchButtonUnique = "set_toggle_telegraph_btn"
)

type TelegraphSwitchButton struct {
	bot  *tb.Bot
	core *core.Core
}

func NewTelegraphSwitchButton(bot *tb.Bot, core *core.Core) *TelegraphSwitchButton {
	return &TelegraphSwitchButton{bot: bot, core: core}
}

func (b *TelegraphSwitchButton) CallbackUnique() string {
	return "\f" + TelegraphSwitchButtonUnique
}

func (b *TelegraphSwitchButton) Description() string {
	return ""
}

func (b *TelegraphSwitchButton) Handle(ctx tb.Context) error {
	c := ctx.Callback()
	if c == nil {
		return ctx.Respond(&tb.CallbackResponse{Text: "error"})
	}

	attachData, err := session.UnmarshalAttachment(ctx.Callback().Data)
	if err != nil {
		return ctx.Respond(&tb.CallbackResponse{Text: "error"})
	}
	subscriberID := attachData.GetUserId()
	if subscriberID != c.Sender.ID {

		channelChat, err := b.bot.ChatByID(subscriberID)
		if err != nil {
			return ctx.Respond(&tb.CallbackResponse{Text: "error"})
		}
		if !chat.IsChatAdmin(b.bot, channelChat, c.Sender.ID) {
			return ctx.Respond(&tb.CallbackResponse{Text: "error"})
		}
	}

	sourceID := uint(attachData.GetSourceId())
	source, _ := b.core.GetSource(context.Background(), sourceID)

	err = b.core.ToggleSubscriptionTelegraph(context.Background(), subscriberID, sourceID)
	if err != nil {
		return ctx.Respond(&tb.CallbackResponse{Text: "error"})
	}

	sub, err := b.core.GetSubscription(context.Background(), subscriberID, sourceID)
	if sub == nil || err != nil {
		return ctx.Respond(&tb.CallbackResponse{Text: "error"})
	}

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

func (b *TelegraphSwitchButton) Middlewares() []tb.MiddlewareFunc {
	return nil
}
