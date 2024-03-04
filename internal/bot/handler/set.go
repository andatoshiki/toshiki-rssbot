package handler

import (
	"bytes"
	"context"
	"fmt"
	"text/template"

	tb "gopkg.in/telebot.v3"

	"github.com/andatoshiki/toshiki-rssbot/internal/bot/chat"
	"github.com/andatoshiki/toshiki-rssbot/internal/bot/session"
	"github.com/andatoshiki/toshiki-rssbot/internal/config"
	"github.com/andatoshiki/toshiki-rssbot/internal/core"
	"github.com/andatoshiki/toshiki-rssbot/internal/model"
)

type Set struct {
	bot  *tb.Bot
	core *core.Core
}

func NewSet(bot *tb.Bot, core *core.Core) *Set {
	return &Set{
		bot:  bot,
		core: core,
	}
}

func (s *Set) Command() string {
	return "/set"
}

func (s *Set) Description() string {
	return "Configure subscriptions"
}

func (s *Set) Handle(ctx tb.Context) error {
	mentionChat, _ := session.GetMentionChatFromCtxStore(ctx)
	ownerID := ctx.Message().Chat.ID
	if mentionChat != nil {
		ownerID = mentionChat.ID
	}

	sources, err := s.core.GetUserSubscribedSources(context.Background(), ownerID)
	if err != nil {
		return ctx.Reply("Failed to fetch subscriptions")
	}
	if len(sources) <= 0 {
		return ctx.Reply("Currently no active subscriptions")
	}

	// Configuration button
	var replyButton []tb.ReplyButton
	replyKeys := [][]tb.ReplyButton{}
	setFeedItemBtns := [][]tb.InlineButton{}
	for _, source := range sources {
		// Add button
		text := fmt.Sprintf("%s %s", source.Title, source.Link)
		replyButton = []tb.ReplyButton{
			{Text: text},
		}
		replyKeys = append(replyKeys, replyButton)
		attachData := &session.Attachment{
			UserId:   ctx.Chat().ID,
			SourceId: uint32(source.ID),
		}

		data := session.Marshal(attachData)
		setFeedItemBtns = append(
			setFeedItemBtns, []tb.InlineButton{
				{
					Unique: SetFeedItemButtonUnique,
					Text:   fmt.Sprintf("[%d] %s", source.ID, source.Title),
					Data:   data,
				},
			},
		)
	}

	return ctx.Reply(
		"Please select the target subscription to configure", &tb.ReplyMarkup{
			InlineKeyboard: setFeedItemBtns,
		},
	)
}

func (s *Set) Middlewares() []tb.MiddlewareFunc {
	return nil
}

const (
	SetFeedItemButtonUnique = "set_feed_item_btn"
	feedSettingTmpl         = `
Subscription<b>Setting</b>
[id] {{ .source.ID }}
[Title] {{ .source.Title }}
[Link] {{.source.Link }}
[Interval] {{if ge .source.ErrorCount .Count }}Pause{{else if lt .source.ErrorCount .Count }}Fetching in progress{{end}}
[Frequency] {{ .sub.Interval }}minute(s)
[Notification] {{if eq .sub.EnableNotification 0}}Disable{{else if eq .sub.EnableNotification 1}}Enable{{end}}
[Telegraph] {{if eq .sub.EnableTelegraph 0}}Disable{{else if eq .sub.EnableTelegraph 1}}Enable{{end}}
[Tag] {{if .sub.Tag}}{{ .sub.Tag }}{{else}}None{{end}}
`
)

type SetFeedItemButton struct {
	bot  *tb.Bot
	core *core.Core
}

func NewSetFeedItemButton(bot *tb.Bot, core *core.Core) *SetFeedItemButton {
	return &SetFeedItemButton{bot: bot, core: core}
}

func (r *SetFeedItemButton) CallbackUnique() string {
	return "\f" + SetFeedItemButtonUnique
}

func (r *SetFeedItemButton) Description() string {
	return ""
}

func (r *SetFeedItemButton) Handle(ctx tb.Context) error {
	attachData, err := session.UnmarshalAttachment(ctx.Callback().Data)
	if err != nil {
		return ctx.Edit("Un")
	}

	subscriberID := attachData.GetUserId()
	// If the subscriber and the button clicker id are not the same, admin permissions need to be verified
	if subscriberID != ctx.Callback().Sender.ID {
		channelChat, err := r.bot.ChatByUsername(fmt.Sprintf("%d", subscriberID))
		if err != nil {
			return ctx.Edit("Failed to fetch subscription information")
		}

		if !chat.IsChatAdmin(r.bot, channelChat, ctx.Callback().Sender.ID) {
			return ctx.Edit("Failed to fetch subscription information") 
		}
	}

	sourceID := uint(attachData.GetSourceId())
	source, err := r.core.GetSource(context.Background(), sourceID)
	if err != nil {
		return ctx.Edit("Unable to find target subscription sources")
	}

	sub, err := r.core.GetSubscription(context.Background(), subscriberID, source.ID)
	if err != nil {
		return ctx.Edit("RSS feed not subscribed by user")
	}

	t := template.New("setting template")
	_, _ = t.Parse(feedSettingTmpl)
	text := new(bytes.Buffer)
	_ = t.Execute(text, map[string]interface{}{"source": source, "sub": sub, "Count": config.ErrorThreshold})
	return ctx.Edit(
		text.String(),
		&tb.SendOptions{ParseMode: tb.ModeHTML},
		&tb.ReplyMarkup{InlineKeyboard: genFeedSetBtn(ctx.Callback(), sub, source)},
	)
}

func genFeedSetBtn(
	c *tb.Callback, sub *model.Subscribe, source *model.Source,
) [][]tb.InlineButton {
	setSubTagKey := tb.InlineButton{
		Unique: SetSubscriptionTagButtonUnique,
		Text:   "Configure tags",
		Data:   c.Data,
	}

	toggleNoticeKey := tb.InlineButton{
		Unique: NotificationSwitchButtonUnique,
		Text:   "Enable notification",
		Data:   c.Data,
	}
	if sub.EnableNotification == 1 {
		toggleNoticeKey.Text = "Disable notification"
	}

	toggleTelegraphKey := tb.InlineButton{
		Unique: TelegraphSwitchButtonUnique,
		Text:   "Enable Telegraph transcoding",
		Data:   c.Data,
	}
	if sub.EnableTelegraph == 1 {
		toggleTelegraphKey.Text = "Disable Telegraph transcoding"
	}

	toggleEnabledKey := tb.InlineButton{
		Unique: SubscriptionSwitchButtonUnique,
		Text:   "Pause update",
		Data:   c.Data,
	}

	if source.ErrorCount >= config.ErrorThreshold {
		toggleEnabledKey.Text = "Restart update"
	}

	feedSettingKeys := [][]tb.InlineButton{
		{
			toggleEnabledKey,
			toggleNoticeKey,
		},
		{
			toggleTelegraphKey,
			setSubTagKey,
		},
	}
	return feedSettingKeys
}

func (r *SetFeedItemButton) Middlewares() []tb.MiddlewareFunc {
	return nil
}
