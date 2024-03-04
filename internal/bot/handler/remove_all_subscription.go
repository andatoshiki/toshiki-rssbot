package handler

import (
	"context"

	tb "gopkg.in/telebot.v3"

	"github.com/andatoshiki/toshiki-rssbot/internal/core"
)

type RemoveAllSubscription struct {
}

func NewRemoveAllSubscription() *RemoveAllSubscription {
	return &RemoveAllSubscription{}
}

func (r RemoveAllSubscription) Command() string {
	return "/unsuball"
}

func (r RemoveAllSubscription) Description() string {
	return "Cancel all existing subscriptions" // Cancel all existing subscriptions
}

func (r RemoveAllSubscription) Handle(ctx tb.Context) error {
	reply := "Unsubscribe all subscription feeds for the current user"
	var confirmKeys [][]tb.InlineButton
	confirmKeys = append(
		confirmKeys, []tb.InlineButton{
			{
				Unique: UnSubAllButtonUnique,
				Text:   "Confirm",
			},
			{
				Unique: CancelUnSubAllButtonUnique,
				Text:   "Cancel",
			},
		},
	)
	return ctx.Reply(reply, &tb.ReplyMarkup{InlineKeyboard: confirmKeys})
}

func (r RemoveAllSubscription) Middlewares() []tb.MiddlewareFunc {
	return nil
}

const (
	UnSubAllButtonUnique       = "unsub_all_confirm_btn"
	CancelUnSubAllButtonUnique = "unsub_all_cancel_btn"
)

type RemoveAllSubscriptionButton struct {
	core *core.Core
}

func NewRemoveAllSubscriptionButton(core *core.Core) *RemoveAllSubscriptionButton {
	return &RemoveAllSubscriptionButton{core: core}
}

func (r *RemoveAllSubscriptionButton) CallbackUnique() string {
	return "\f" + UnSubAllButtonUnique
}

func (r *RemoveAllSubscriptionButton) Description() string {
	return ""
}

func (r *RemoveAllSubscriptionButton) Handle(ctx tb.Context) error {
	err := r.core.UnsubscribeAllSource(context.Background(), ctx.Sender().ID)
	if err != nil {
		return ctx.Edit("Failed to unsubscribe")
	}
	return ctx.Edit("Successfully unsubscribed")
}

func (r *RemoveAllSubscriptionButton) Middlewares() []tb.MiddlewareFunc {
	return nil
}

type CancelRemoveAllSubscriptionButton struct {
}

func NewCancelRemoveAllSubscriptionButton() *CancelRemoveAllSubscriptionButton {
	return &CancelRemoveAllSubscriptionButton{}
}

func (r *CancelRemoveAllSubscriptionButton) CallbackUnique() string {
	return "\f" + CancelUnSubAllButtonUnique
}

func (r *CancelRemoveAllSubscriptionButton) Description() string {
	return ""
}

func (r *CancelRemoveAllSubscriptionButton) Handle(ctx tb.Context) error {
	return ctx.Edit("Cancel current operations")
}

func (r *CancelRemoveAllSubscriptionButton) Middlewares() []tb.MiddlewareFunc {
	return nil
}
