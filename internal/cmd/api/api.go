package api

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	bolt "go.etcd.io/bbolt"

	"github.com/Anton-Kraev/gopay"
	"github.com/Anton-Kraev/gopay/internal/client/minio"
	"github.com/Anton-Kraev/gopay/internal/client/yookassa"
	"github.com/Anton-Kraev/gopay/internal/http/handler"
	"github.com/Anton-Kraev/gopay/internal/http/server"
	"github.com/Anton-Kraev/gopay/internal/links"
	"github.com/Anton-Kraev/gopay/internal/logger"
	repo "github.com/Anton-Kraev/gopay/internal/repository/bolt"
	"github.com/Anton-Kraev/gopay/internal/validator"
)

type API struct {
	Env                 string
	GopayHost           string
	GopayPort           string
	DBFilePath          string
	DBOpenTimeout       time.Duration
	YookassaCheckoutURL string
	YookassaShopID      string
	YookassaAPIToken    string
	MinioBucketName     string
	MinioURL            string
	MinioUser           string
	MinioPassword       string
}

func (a *API) Start(ctx context.Context) error {
	log := logger.Setup(a.Env)
	log.Info("Config parsed", slog.Any("config", a))

	db, err := bolt.Open(
		a.DBFilePath,
		0600,
		&bolt.Options{Timeout: a.DBOpenTimeout},
	)
	if err != nil {
		return err
	}

	defer func(db *bolt.DB) {
		if err = db.Close(); err != nil {
			log.Error(err.Error())
		}
	}(db)

	paymentStorage, err := repo.NewPaymentRepository(db)
	if err != nil {
		return err
	}

	paymentService := yookassa.NewClient(yookassa.Config{
		CheckoutURL: a.YookassaCheckoutURL,
		ShopID:      a.YookassaShopID,
		APIToken:    a.YookassaAPIToken,
	})

	linkGenerator := links.NewGenerator(fmt.Sprintf("%s:%s", a.GopayHost, a.GopayPort))

	pm := gopay.NewPaymentManager(
		linkGenerator,
		paymentStorage,
		paymentService,
	)

	fileStorage, err := minio.NewClient(ctx, minio.Config{
		BucketName: a.MinioBucketName,
		URL:        a.MinioURL,
		User:       a.MinioUser,
		Password:   a.MinioPassword,
	})
	if err != nil {
		return err
	}

	hndl := handler.NewHandler(pm, fileStorage)

	val, err := validator.NewValidator()
	if err != nil {
		return err
	}

	srv := server.NewServer(hndl, log, val)
	echoSrv := srv.InitRoutes()

	return echoSrv.Start(":" + a.GopayPort)
}
