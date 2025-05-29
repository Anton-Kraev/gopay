package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/Anton-Kraev/gopay/internal/cmd/api"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	cmd := api.NewAPICmd()
	if err := cmd.Run(ctx, os.Args); err != nil {
		log.Fatalln(err)
	}
}
