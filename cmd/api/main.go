package main

import (
	"fmt"
	"log"

	bolt "go.etcd.io/bbolt"

	"github.com/Anton-Kraev/gopay"
	"github.com/Anton-Kraev/gopay/internal/client/yookassa"
	"github.com/Anton-Kraev/gopay/internal/config/api"
	"github.com/Anton-Kraev/gopay/internal/http/handler"
	"github.com/Anton-Kraev/gopay/internal/http/server"
	"github.com/Anton-Kraev/gopay/internal/links"
	"github.com/Anton-Kraev/gopay/internal/logger"
	repo "github.com/Anton-Kraev/gopay/internal/repository/bolt"
	"github.com/Anton-Kraev/gopay/internal/templates"
	"github.com/Anton-Kraev/gopay/internal/validator"
	"github.com/Anton-Kraev/gopay/mock"
)

const (
	configPath = "./configs/api.yaml"
)

func main() {
	cfg, err := api.GetConfig(configPath)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(cfg)

	db, err := bolt.Open(
		cfg.DB.FilePath,
		0600,
		&bolt.Options{Timeout: cfg.DB.OpenTimeout},
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

	paymentService := yookassa.NewClient(cfg.Yookassa.CheckoutURL, yookassa.AuthConfig{
		ID:    cfg.Yookassa.ShopID,
		Token: cfg.Yookassa.APIToken,
	})

	linkGenerator := links.NewGenerator(fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port))

	templateStorage, err := templates.New(db)
	if err != nil {
		log.Fatalln(err)
	}

	if err = templateStorage.SetTemplate("template", gopay.PaymentTemplate{
		Currency:     "RUB",
		Amount:       1,
		Description:  "description",
		ResourceLink: "http://127.0.0.1:8080/api/files/123",
	}); err != nil {
		log.Fatalln(err)
	}

	pm := gopay.NewPaymentManager(
		templateStorage,
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
	echoSrv.Logger.Fatal(echoSrv.Start(":" + cfg.Server.Port))
}
