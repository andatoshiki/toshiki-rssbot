package handler

import tb "gopkg.in/telebot.v3"

type Import struct {
}

func NewImport() *Import {
	return &Import{}
}

func (i *Import) Command() string {
	return "/import"
}

func (i *Import) Description() string {
	return "Import OPML file"
}

func (i *Import) Handle(ctx tb.Context) error {
	reply := "Please append and send your OPML attachment as a direct message to the bot; if you need to import OPML feed on behalf of your channel, please affix the channel ID with the attachment; eg: @toshikidev"
	return ctx.Reply(reply)
}

func (i *Import) Middlewares() []tb.MiddlewareFunc {
	return nil
}
