package api

import (
	"context"
	"fmt"
	"time"

	"github.com/urfave/cli/v3"
)

func NewAPICmd() *cli.Command {
	var api API

	cmd := &cli.Command{
		Name:        "api",
		Usage:       "Run GoPay API",
		Description: "GoPay API",
		UsageText: "api " +
			"--yookassa-checkout-url <gopay_checkout> --yookassa-shop-id <shop_id> --yookassa-api-token <api_token> " +
			"--minio-user <user> --minio-password <password>",
		Action: func(ctx context.Context, _ *cli.Command) error {
			if err := api.Start(ctx); err != nil {
				return fmt.Errorf("Api.Start: %w", err)
			}

			return nil
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "env",
				Usage:       "Environment type (dev/prod)",
				Value:       "dev",
				Sources:     cli.EnvVars("ENV"),
				Destination: &api.Env,
			},
			&cli.StringFlag{
				Name:        "gopay-host",
				Usage:       "GoPay server host",
				Value:       "localhost",
				Sources:     cli.EnvVars("GOPAY_HOST"),
				Destination: &api.GopayHost,
			},
			&cli.StringFlag{
				Name:        "gopay-port",
				Aliases:     []string{"p"},
				Usage:       "GoPay server port",
				Value:       "8080",
				Sources:     cli.EnvVars("GOPAY_PORT"),
				Destination: &api.GopayPort,
			},
			&cli.StringFlag{
				Name:        "db-file-path",
				Usage:       "Database file path",
				Value:       "data.db",
				Sources:     cli.EnvVars("DB_FILE"),
				Destination: &api.DBFilePath,
			},
			&cli.DurationFlag{
				Name:        "db-open-timeout",
				Usage:       "Database open timeout",
				Value:       10 * time.Second,
				Sources:     cli.EnvVars("DB_OPEN_TIMEOUT"),
				Destination: &api.DBOpenTimeout,
			},
			&cli.StringFlag{
				Name:        "yookassa-checkout-url",
				Usage:       "Yookassa checkout URL",
				Required:    true,
				Sources:     cli.EnvVars("YOOKASSA_CHECKOUT_URL"),
				Destination: &api.YookassaCheckoutURL,
			},
			&cli.StringFlag{
				Name:        "yookassa-shop-id",
				Usage:       "Yookassa Shop ID",
				Required:    true,
				Sources:     cli.EnvVars("YOOKASSA_SHOP_ID"),
				Destination: &api.YookassaShopID,
			},
			&cli.StringFlag{
				Name:        "yookassa-api-token",
				Usage:       "Yookassa API token",
				Required:    true,
				Sources:     cli.EnvVars("YOOKASSA_API_TOKEN"),
				Destination: &api.YookassaAPIToken,
			},
			&cli.StringFlag{
				Name:        "minio-bucket-name",
				Usage:       "MinIO bucket name",
				Value:       "geopdfs",
				Sources:     cli.EnvVars("MINIO_BUCKET_NAME"),
				Destination: &api.MinioBucketName,
			},
			&cli.StringFlag{
				Name:        "minio-url",
				Usage:       "MinIO URL",
				Value:       "http://127.0.0.1:9000",
				Sources:     cli.EnvVars("MINIO_URL"),
				Destination: &api.MinioURL,
			},
			&cli.StringFlag{
				Name:        "minio-user",
				Usage:       "MinIO user",
				Required:    true,
				Sources:     cli.EnvVars("MINIO_USER"),
				Destination: &api.MinioUser,
			},
			&cli.StringFlag{
				Name:        "minio-password",
				Usage:       "MinIO password",
				Required:    true,
				Sources:     cli.EnvVars("MINIO_PASSWORD"),
				Destination: &api.MinioPassword,
			},
		},
	}

	return cmd
}
