package server

import (
	"log/slog"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	slogecho "github.com/samber/slog-echo"
	swagecho "github.com/swaggo/echo-swagger"

	// Register generated Swagger docs
	_ "github.com/Anton-Kraev/gopay/docs"
	"github.com/Anton-Kraev/gopay/internal/validator"
)

type handlers interface {
	NewPayment(c echo.Context) error
	AllPayment(c echo.Context) error
	GetPayment(c echo.Context) error
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

// InitRoutes init API routes and middlewares
// @title GoPay API
// @version 1.0
// @description API for payment processing and digital goods access management
// @license.name MIT license
// @license.url https://opensource.org/licenses/MIT
// @contact.name Author's contact
// @contact.url https://t.me/iksvayai
// @BasePath /api
func (s Server) InitRoutes() *echo.Echo {
	e := echo.New()

	e.Validator = s.validator

	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	e.Use(slogecho.New(s.logger))

	e.GET("/swagger/*", swagecho.WrapHandler)

	g := e.Group("/api")

	g.POST("/payments", s.handlers.NewPayment)
	g.GET("/payments", s.handlers.AllPayment)
	g.GET("/payments/:id", s.handlers.GetPayment)
	g.GET("/:id", s.handlers.Redirect)
	g.POST("/checkout", s.handlers.Checkout)
	g.GET("/files/:id", s.handlers.File)

	return e
}
