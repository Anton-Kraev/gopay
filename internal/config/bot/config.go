package bot

import (
	"context"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"
)

type Config struct {
	GopayServerURL string
	TGBotToken     string
	TGAdminIDs     string
}

func LoadConfig(ctx context.Context) (Config, error) {
	cmd := &cli.Command{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "gopay-server-url",
				Usage:   "GoPay server url",
				Value:   "http://127.0.0.1:8080",
				Sources: cli.EnvVars("GOPAY_SERVER_URL"),
			},
			&cli.StringFlag{
				Name:     "tg-bot-token",
				Usage:    "Token for Telegram bot API",
				Required: true,
				Sources:  cli.EnvVars("TG_BOT_TOKEN"),
			},
			&cli.StringFlag{
				Name:     "tg-admin-ids",
				Usage:    "Admins Telegram identifiers in format <id1>,<id2>,<id3>",
				Required: true,
				Sources:  cli.EnvVars("TG_ADMIN_IDS"),
			},
		},
	}

	if err := cmd.Run(ctx, os.Args); err != nil {
		return Config{}, fmt.Errorf("bot.LoadConfig: %w", err)
	}

	cfg := Config{
		GopayServerURL: cmd.String("gopay-server-url"),
		TGBotToken:     cmd.String("tg-bot-token"),
		TGAdminIDs:     cmd.String("tg-admin-ids"),
	}

	return cfg, nil
}
