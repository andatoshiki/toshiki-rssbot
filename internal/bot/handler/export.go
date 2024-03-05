package handler

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"go.uber.org/zap"
	tb "gopkg.in/telebot.v3"

	"github.com/andatoshiki/toshiki-rssbot/internal/bot/message"
	"github.com/andatoshiki/toshiki-rssbot/internal/core"
	"github.com/andatoshiki/toshiki-rssbot/internal/log"
	"github.com/andatoshiki/toshiki-rssbot/internal/model"
	"github.com/andatoshiki/toshiki-rssbot/internal/opml"
)

type Export struct {
	core *core.Core
}

func NewExport(core *core.Core) *Export {
	return &Export{core: core}
}

func (e *Export) Description() string {
	return "Export all subscribed feeds to OML"
}

func (e *Export) Command() string {
	return "/export"
}

func (e *Export) getChannelSources(bot *tb.Bot, opUserID int64, channelName string) ([]*model.Source, error) {
	// export channel subscription sources
	channelChat, err := bot.ChatByUsername(channelName)
	if err != nil {
		return nil, errors.New("Failed to fetch channel information from remote origin}")
	}

	adminList, err := bot.AdminsOf(channelChat)
	if err != nil {
		return nil, errors.New("Failed to fetch information of channel administrator")
	}

	senderIsAdmin := false
	for _, admin := range adminList {
		if opUserID == admin.User.ID {
			senderIsAdmin = true
			break
		}
	}

	if !senderIsAdmin {
		return nil, errors.New("Bot operational executions by non-administrative users are not permitted")
	}

	sources, err := e.core.GetUserSubscribedSources(context.Background(), channelChat.ID)
	if err != nil {
		zap.S().Error(err)
		return nil, errors.New("Failed to fetch resources or information from remote source origin")
	}
	return sources, nil
}

func (e *Export) Handle(ctx tb.Context) error {
	mention := message.MentionFromMessage(ctx.Message())
	var sources []*model.Source
	if mention == "" {
		var err error
		sources, err = e.core.GetUserSubscribedSources(context.Background(), ctx.Chat().ID)
		if err != nil {
			log.Error(err)
			return ctx.Send("Failed to export") //
		}
	} else {
		var err error
		sources, err = e.getChannelSources(ctx.Bot(), ctx.Chat().ID, mention)
		if err != nil {
			log.Error(err)
			return ctx.Send("Failed to export")
		}
	}

	if len(sources) == 0 {
		return ctx.Send("The subscription list is currently empty ")
	}

	opmlStr, err := opml.ToOPML(sources)
	if err != nil {
		return ctx.Send("Failed to export")
	}
	opmlFile := &tb.Document{File: tb.FromReader(strings.NewReader(opmlStr))}
	opmlFile.FileName = fmt.Sprintf("subscriptions_%d.opml", time.Now().Unix())
	if err := ctx.Send(opmlFile); err != nil {
		log.Errorf("send OPML file failed, err:%v", err)
		return ctx.Send("Failed to export")
	}
	return nil
}

func (e *Export) Middlewares() []tb.MiddlewareFunc {
	return nil
}
