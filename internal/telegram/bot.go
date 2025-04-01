package telegram

import (
	"context"
	"fmt"

	"github.com/mymmrac/telego"
)

type Telegram struct {
	bot       *telego.Bot
	fsm       map[int64]string // finite state machine (store user states in format "user_id: state")
	whitelist []int64          // list of allowed users (gopay admins user ids)
}

func New(token string, adminIDs []int64) (*Telegram, error) {
	tgBot, err := telego.NewBot(token)
	if err != nil {
		return nil, fmt.Errorf("telegram.New: %w", err)
	}

	return &Telegram{
		bot:       tgBot,
		fsm:       make(map[int64]string),
		whitelist: adminIDs,
	}, nil
}

func (t *Telegram) Start(ctx context.Context) error {
	updates, err := t.bot.UpdatesViaLongPolling(ctx, nil)
	if err != nil {
		return fmt.Errorf("telegram.Telegram.Start: %w", err)
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case update := <-updates:
			t.handleUpdate(ctx, update)
		}
	}
}

func (t *Telegram) handleUpdate(ctx context.Context, update telego.Update) {
	_, _ = t.bot.SendMessage(ctx, &telego.SendMessageParams{
		ChatID: telego.ChatID{ID: update.Message.Chat.ID},
		Text:   update.Message.Text,
	})
}
