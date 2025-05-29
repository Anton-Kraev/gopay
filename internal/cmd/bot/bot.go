package bot

import (
	"context"
	"log/slog"

	"github.com/Anton-Kraev/gopay"
	"github.com/Anton-Kraev/gopay/internal/logger"
	"github.com/Anton-Kraev/gopay/internal/telegram"
	"github.com/Anton-Kraev/gopay/internal/typeconv"
)

type Bot struct {
	Env            string
	GopayServerURL string
	TGBotToken     string
	TGAdminIDs     string
}

func (b *Bot) Start(ctx context.Context) error {
	log := logger.Setup(b.Env)
	log.Info("Config parsed", slog.Any("config", b))

	adminIDs, err := typeconv.StringToInt64Slice(b.TGAdminIDs)
	if err != nil {
		return err
	}

	adminClient, err := gopay.NewAdminClient(b.GopayServerURL)
	if err != nil {
		return err
	}

	tg, err := telegram.New(telegram.Config{
		BotToken: b.TGBotToken,
		AdminIDs: adminIDs,
	}, adminClient, logger.Setup(b.Env))
	if err != nil {
		return err
	}

	log.Info("Starting tg-bot")

	return tg.Start(ctx)
}
