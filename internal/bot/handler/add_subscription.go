package handler

import (
	"context"
	"errors"
	"fmt"

	"go.uber.org/zap"
	tb "gopkg.in/telebot.v3"

	"github.com/andatoshiki/toshiki-rssbot/internal/bot/message"
	"github.com/andatoshiki/toshiki-rssbot/internal/core"
	"github.com/andatoshiki/toshiki-rssbot/internal/log"
)

type AddSubscription struct {
	core *core.Core
}

func NewAddSubscription(core *core.Core) *AddSubscription {
	return &AddSubscription{
		core: core,
	}
}

func (a *AddSubscription) Command() string {
	return "/sub"
}

func (a *AddSubscription) Description() string {
	return "Subscribe RSS feed source"
}

func (a *AddSubscription) addSubscriptionForChat(ctx tb.Context) error {
	sourceURL := message.URLFromMessage(ctx.Message())
	if sourceURL == "" {
		// 未附带链接，使用
		hint := fmt.Sprintf("Please append the target subscription url for RSS feed at the end of the command; e.g.: %s https://github.blog/feed/", a.Command()) // Please append the target subscription url for RSS feed at the end of the command; e.g.: %s https://github.blog/feed/
		return ctx.Send(hint, &tb.SendOptions{ReplyTo: ctx.Message()})
	}

	source, err := a.core.CreateSource(context.Background(), sourceURL)
	if err != nil {
		return ctx.Reply(fmt.Sprintf("%s, failed to subscribe", err))
	}

	log.Infof("%d subscribe [%d]%s %s", ctx.Chat().ID, source.ID, source.Title, source.Link)
	if err := a.core.AddSubscription(context.Background(), ctx.Chat().ID, source.ID); err != nil {
		if err == core.ErrSubscriptionExist {
			return ctx.Reply("Source subscribed and exist in present feed list already, please do not repeatedly duplicate subscription")
		}
		log.Errorf("add subscription user %d source %d failed %v", ctx.Chat().ID, source.ID, err)
		return ctx.Reply("Failed to subscribe from source")
	}

	return ctx.Reply(
		fmt.Sprintf("[[%d]][%s](%s) Successfully subscribed from source to feed list", source.ID, source.Title, source.Link),
		&tb.SendOptions{
			DisableWebPagePreview: true,
			ParseMode:             tb.ModeMarkdown,
		},
	)
}

func (a *AddSubscription) hasChannelPrivilege(bot *tb.Bot, channelChat *tb.Chat, opUserID int64, botID int64) (
	bool, error,
) {
	adminList, err := bot.AdminsOf(channelChat)
	if err != nil {
		zap.S().Error(err)
		return false, errors.New("Failed to fetch channel information from remote origin")
	}

	senderIsAdmin := false
	botIsAdmin := false
	for _, admin := range adminList {
		if opUserID == admin.User.ID {
			senderIsAdmin = true
		}
		if botID == admin.User.ID {
			botIsAdmin = true
		}
	}

	return botIsAdmin && senderIsAdmin, nil
}

func (a *AddSubscription) addSubscriptionForChannel(ctx tb.Context, channelName string) error {
	sourceURL := message.URLFromMessage(ctx.Message())
	if sourceURL == "" {
		return ctx.Send("Please run ' /sub `@channel_id` URL ' command for subscription on behalf of a specific channel; e.g.: @toshikidev", &tb.SendOptions{ReplyTo: ctx.Message()})
	}

	bot := ctx.Bot()
	channelChat, err := bot.ChatByUsername(channelName)
	if err != nil {
		return ctx.Reply("Failed to fetch channel details & information")
	}
	if channelChat.Type != tb.ChatChannel {
		return ctx.Reply("Either you or the bot is currently not the administrator of the channel provided, failed to configure subscription")
	}

	hasPrivilege, err := a.hasChannelPrivilege(bot, channelChat, ctx.Sender().ID, bot.Me.ID)
	if err != nil {
		return ctx.Reply(err.Error())
	}
	if !hasPrivilege {
		return ctx.Reply("Either you or the bot is currently not the administrator of the channel provided, failed to configure subscription")
	}

	source, err := a.core.CreateSource(context.Background(), sourceURL)
	if err != nil {
		return ctx.Reply(fmt.Sprintf("%s，Subscription failed", err))
	}

	log.Infof("%d subscribe [%d]%s %s", channelChat.ID, source.ID, source.Title, source.Link)
	if err := a.core.AddSubscription(context.Background(), channelChat.ID, source.ID); err != nil {
		if err == core.ErrSubscriptionExist {
			return ctx.Reply("The source has been subscribed, please do not repeatedly duplicate subscription")
		}
		log.Errorf("add subscription user %d source %d failed %v", channelChat.ID, source.ID, err)
		return ctx.Reply("Failed to subscribe from source")
	}

	return ctx.Reply(
		fmt.Sprintf("[[%d]] [%s](%s) Feed source has successfully subscribed", source.ID, source.Title, source.Link),
		&tb.SendOptions{
			DisableWebPagePreview: true,
			ParseMode:             tb.ModeMarkdown,
		},
	)
}

func (a *AddSubscription) Handle(ctx tb.Context) error {
	mention := message.MentionFromMessage(ctx.Message())
	if mention != "" {
		// has mention, add subscription for channel
		return a.addSubscriptionForChannel(ctx, mention)
	}
	return a.addSubscriptionForChat(ctx)
}

func (a *AddSubscription) Middlewares() []tb.MiddlewareFunc {
	return nil
}
