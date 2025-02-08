package main

import (
	"github.com/Anton-Kraev/gopay"
	"github.com/Anton-Kraev/gopay/internal/http/handler"
	"github.com/Anton-Kraev/gopay/internal/http/server"
	"github.com/Anton-Kraev/gopay/internal/logger"
	"github.com/Anton-Kraev/gopay/internal/validator"
	"github.com/Anton-Kraev/gopay/mock"
)

func main() {
	pm := gopay.NewPaymentManager(
		mock.NewTemplates(),
		mock.NewLinkGenerator(),
		mock.NewPaymentStorage(),
		mock.NewPaymentService(),
	)

	fileStorage := mock.FileStorage{}

	hndl := handler.NewHandler(pm, nil, fileStorage)

	log := logger.Setup("local")
	val, err := validator.NewValidator()
	if err != nil {
		panic(err)
	}

	srv := server.NewServer(hndl, log, val)

	echoSrv := srv.InitRoutes()
	echoSrv.Logger.Fatal(echoSrv.Start(":1323"))
}
