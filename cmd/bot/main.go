package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/Anton-Kraev/gopay"
	"github.com/Anton-Kraev/gopay/internal/config/bot"
	"github.com/Anton-Kraev/gopay/internal/logger"
	"github.com/Anton-Kraev/gopay/internal/telegram"
	"github.com/Anton-Kraev/gopay/internal/typeconv"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	cfg, err := bot.GetConfig()
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(cfg)

	adminIDs, err := typeconv.StringToInt64Slice(cfg.AdminIDs)
	if err != nil {
		log.Fatalln(err)
	}

	adminClient := gopay.NewAdminClient(cfg.ServerURL)

	tg, err := telegram.New(adminClient, cfg.Token, adminIDs, logger.Setup("local"))
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Starting tg-bot")
	log.Fatalln(tg.Start(ctx))
}
