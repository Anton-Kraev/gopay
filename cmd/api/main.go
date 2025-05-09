package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	bolt "go.etcd.io/bbolt"

	"github.com/Anton-Kraev/gopay"
	"github.com/Anton-Kraev/gopay/internal/client/yookassa"
	"github.com/Anton-Kraev/gopay/internal/config/api"
	"github.com/Anton-Kraev/gopay/internal/http/handler"
	"github.com/Anton-Kraev/gopay/internal/http/server"
	"github.com/Anton-Kraev/gopay/internal/links"
	"github.com/Anton-Kraev/gopay/internal/logger"
	repo "github.com/Anton-Kraev/gopay/internal/repository/bolt"
	"github.com/Anton-Kraev/gopay/internal/validator"
	"github.com/Anton-Kraev/gopay/mock"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	cfg, err := api.LoadConfig(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(cfg)

	db, err := bolt.Open(
		cfg.DBFilePath,
		0600,
		&bolt.Options{Timeout: cfg.DBOpenTimeout},
	)
	if err != nil {
		log.Fatalln(err)
	}

	defer func(db *bolt.DB) {
		if err = db.Close(); err != nil {
			log.Println("error closing db connection:", err)
		}
	}(db)

	paymentStorage, err := repo.NewPaymentRepository(db)
	if err != nil {
		log.Fatalln(err)
	}

	paymentService := yookassa.NewClient(yookassa.Config{
		CheckoutURL: cfg.YookassaCheckoutURL,
		ShopID:      cfg.YookassaShopID,
		APIToken:    cfg.YookassaAPIToken,
	})

	linkGenerator := links.NewGenerator(fmt.Sprintf("%s:%s", cfg.GopayHost, cfg.GopayPort))

	pm := gopay.NewPaymentManager(
		linkGenerator,
		paymentStorage,
		paymentService,
	)

	fileStorage := mock.FileStorage{}

	hndl := handler.NewHandler(pm, fileStorage)

	val, err := validator.NewValidator()
	if err != nil {
		log.Fatalln(err)
	}

	srv := server.NewServer(hndl, logger.Setup(cfg.Env), val)

	echoSrv := srv.InitRoutes()
	echoSrv.Logger.Fatal(echoSrv.Start(":" + cfg.GopayPort))
}
