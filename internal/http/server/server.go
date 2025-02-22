package server

import (
	"log/slog"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	slogecho "github.com/samber/slog-echo"

	"github.com/Anton-Kraev/gopay/internal/validator"
)

type handlers interface {
	NewPayment(c echo.Context) error
	Redirect(c echo.Context) error
	Checkout(c echo.Context) error
	File(c echo.Context) error
}

type Server struct {
	handlers  handlers
	logger    *slog.Logger
	validator *validator.Validator
}

func NewServer(handlers handlers, logger *slog.Logger, validator *validator.Validator) Server {
	return Server{
		handlers:  handlers,
		logger:    logger,
		validator: validator,
	}
}

func (s Server) InitRoutes() *echo.Echo {
	e := echo.New()

	e.Validator = s.validator

	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	e.Use(slogecho.New(s.logger))

	g := e.Group("/api")

	g.POST("/:id", s.handlers.NewPayment)
	g.GET("/:id", s.handlers.Redirect)
	g.POST("/payments", s.handlers.Checkout)
	g.GET("/files/:id", s.handlers.File)

	return e
}
