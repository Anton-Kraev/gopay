package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/Anton-Kraev/gopay/internal/cmd/bot"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	cmd := bot.NewBotCmd()
	if err := cmd.Run(ctx, os.Args); err != nil {
		log.Fatalln(err)
	}
}
