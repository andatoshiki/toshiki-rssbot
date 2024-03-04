package middleware

import (
	"github.com/andatoshiki/toshiki-rssbot/internal/bot/chat"
	"github.com/andatoshiki/toshiki-rssbot/internal/bot/session"

	tb "gopkg.in/telebot.v3"
)

func IsChatAdmin() tb.MiddlewareFunc {
	return func(next tb.HandlerFunc) tb.HandlerFunc {
		return func(c tb.Context) error {
			if !chat.IsChatAdmin(c.Bot(), c.Chat(), c.Sender().ID) {
				return c.Reply("You are not the administrator of the current session")
			}

			v := c.Get(session.StoreKeyMentionChat.String())
			if v != nil {
				mentionChat, ok := v.(*tb.Chat)
				if !ok {
					return c.Reply("Internal service error")
				}
				if !chat.IsChatAdmin(c.Bot(), mentionChat, c.Sender().ID) {
					return c.Reply("You are not the administrator of the current session")
				}
			}
			return next(c)
		}
	}
}
