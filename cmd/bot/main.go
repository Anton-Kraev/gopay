package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/Anton-Kraev/gopay/internal/config/bot"
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

	tg, err := telegram.New(cfg.Token, adminIDs)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Starting tg-bot")
	log.Fatalln(tg.Start(ctx))
}
