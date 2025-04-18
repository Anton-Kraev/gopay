package telegram

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/mymmrac/telego"

	"github.com/Anton-Kraev/gopay"
)

type Telegram struct {
	adminClient       gopay.AdminClient
	bot               *telego.Bot
	fsm               map[int64]state                   // finite state machine (store user states in format "user_id: state")
	newPaymentService map[int64]gopay.NewPaymentService // service with user data for creating new payments
	whitelist         []int64                           // list of allowed users (gopay admins user ids)
	log               *slog.Logger
}

func New(adminClient gopay.AdminClient, token string, adminIDs []int64, log *slog.Logger) (*Telegram, error) {
	tgBot, err := telego.NewBot(token)
	if err != nil {
		return nil, fmt.Errorf("telegram.New: %w", err)
	}

	return &Telegram{
		adminClient:       adminClient,
		bot:               tgBot,
		fsm:               make(map[int64]state),
		newPaymentService: make(map[int64]gopay.NewPaymentService),
		whitelist:         adminIDs,
		log:               log,
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
