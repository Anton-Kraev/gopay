package api

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/urfave/cli/v3"
)

type Config struct {
	Env                 string
	GopayHost           string
	GopayPort           string
	DBFilePath          string
	DBOpenTimeout       time.Duration
	YookassaCheckoutURL string
	YookassaShopID      string
	YookassaAPIToken    string
}

func LoadConfig(ctx context.Context) (Config, error) {
	cmd := &cli.Command{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "env",
				Usage:   "Environment type (dev/prod)",
				Value:   "dev",
				Sources: cli.EnvVars("ENV"),
			},
			&cli.StringFlag{
				Name:    "gopay-host",
				Usage:   "GoPay server host",
				Value:   "localhost",
				Sources: cli.EnvVars("GOPAY_HOST"),
			},
			&cli.StringFlag{
				Name:    "gopay-port",
				Aliases: []string{"p"},
				Usage:   "GoPay server port",
				Value:   "8080",
				Sources: cli.EnvVars("GOPAY_PORT"),
			},
			&cli.StringFlag{
				Name:    "db-file-path",
				Usage:   "Database file path",
				Value:   "data.db",
				Sources: cli.EnvVars("DB_FILE"),
			},
			&cli.DurationFlag{
				Name:    "db-open-timeout",
				Usage:   "Database open timeout",
				Value:   10 * time.Second,
				Sources: cli.EnvVars("DB_OPEN_TIMEOUT"),
			},
			&cli.StringFlag{
				Name:     "yookassa-checkout-url",
				Usage:    "Yookassa checkout URL",
				Required: true,
				Sources:  cli.EnvVars("YOOKASSA_CHECKOUT_URL"),
			},
			&cli.StringFlag{
				Name:     "yookassa-shop-id",
				Usage:    "Yookassa Shop ID",
				Required: true,
				Sources:  cli.EnvVars("YOOKASSA_SHOP_ID"),
			},
			&cli.StringFlag{
				Name:     "yookassa-api-token",
				Usage:    "Yookassa API token",
				Required: true,
				Sources:  cli.EnvVars("YOOKASSA_API_TOKEN"),
			},
		},
	}

	if err := cmd.Run(ctx, os.Args); err != nil {
		return Config{}, fmt.Errorf("api.LoadConfig: %w", err)
	}

	cfg := Config{
		Env:                 cmd.String("env"),
		GopayHost:           cmd.String("gopay-host"),
		GopayPort:           cmd.String("gopay-port"),
		DBFilePath:          cmd.String("db-file-path"),
		DBOpenTimeout:       cmd.Duration("db-open-timeout"),
		YookassaCheckoutURL: cmd.String("yookassa-checkout-url"),
		YookassaShopID:      cmd.String("yookassa-shop-id"),
		YookassaAPIToken:    cmd.String("yookassa-api-token"),
	}

	return cfg, nil
}
