package handler

import (
	"fmt"

	tb "gopkg.in/telebot.v3"

	"github.com/andatoshiki/toshiki-rssbot/internal/log"
)

type Start struct {
}

func NewStart() *Start {
	return &Start{}
}

func (s *Start) Command() string {
	return "/start"
}

func (s *Start) Description() string {
	return "Start using the bot"
}

func (s *Start) Handle(ctx tb.Context) error {
	log.Infof("/start id: %d", ctx.Chat().ID)
	return ctx.Send(fmt.Sprintf("H"))
}

func (s *Start) Middlewares() []tb.MiddlewareFunc {
	return nil
}
