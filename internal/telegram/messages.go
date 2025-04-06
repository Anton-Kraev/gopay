package telegram

import (
	"context"
	"fmt"

	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

func (t *Telegram) sendMessage(ctx context.Context, update telego.Update, handler, msg string) error {
	_, err := t.bot.SendMessage(ctx, tu.Message(
		tu.ID(update.Message.Chat.ID),
		msg,
	))
	if err != nil {
		return fmt.Errorf("%s: %w", handler, err)
	}

	return nil
}
