package bot

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
)

func NewBotCmd() *cli.Command {
	var bot Bot

	cmd := &cli.Command{
		Name:        "bot",
		Usage:       "Run GoPay Telegram Bot",
		Description: "GoPay Telegram Bot",
		UsageText:   "bot --tg-bot-token <token> --tg-admin-ids <id1>,<id2>",
		Action: func(ctx context.Context, _ *cli.Command) error {
			if err := bot.Start(ctx); err != nil {
				return fmt.Errorf("Bot.Start: %w", err)
			}

			return nil
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "env",
				Usage:       "Environment type (dev/prod)",
				Value:       "dev",
				Sources:     cli.EnvVars("ENV"),
				Destination: &bot.Env,
			},
			&cli.StringFlag{
				Name:        "gopay-server-url",
				Usage:       "GoPay server url",
				Value:       "http://127.0.0.1:8080",
				Sources:     cli.EnvVars("GOPAY_SERVER_URL"),
				Destination: &bot.GopayServerURL,
			},
			&cli.StringFlag{
				Name:        "tg-bot-token",
				Usage:       "Token for Telegram bot API",
				Required:    true,
				Sources:     cli.EnvVars("TG_BOT_TOKEN"),
				Destination: &bot.TGBotToken,
			},
			&cli.StringFlag{
				Name:        "tg-admin-ids",
				Usage:       "Admins Telegram identifiers in format <id1>,<id2>,<id3>",
				Required:    true,
				Sources:     cli.EnvVars("TG_ADMIN_IDS"),
				Destination: &bot.TGAdminIDs,
			},
		},
	}

	return cmd
}
