package server

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	slogecho "github.com/samber/slog-echo"

	"github.com/Anton-Kraev/gopay/internal/logger"
)

type handlers interface {
	NewPayment(c echo.Context) error
	Redirect(c echo.Context) error
	Checkout(c echo.Context) error
	File(c echo.Context) error
}

func InitRoutes(env string, handlers handlers) *echo.Echo {
	e := echo.New()

	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	e.Use(slogecho.New(logger.Setup(env)))

	e.POST("/:id", handlers.NewPayment)
	e.GET("/:id", handlers.Redirect)
	e.POST("/payment/:id", handlers.Checkout)
	e.GET("/file/:id", handlers.File)

	return e
}
