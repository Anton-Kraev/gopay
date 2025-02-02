package handler

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/Anton-Kraev/gopay"
)

type (
	auth        interface{}
	fileStorage interface {
		GetData(id gopay.ID) ([]byte, error)
	}
)

type Handler struct {
	paymentManager *gopay.PaymentManager
	auth           auth
	fileStorage    fileStorage
}

func NewHandler(paymentManager *gopay.PaymentManager, auth auth, fileStorage fileStorage) Handler {
	return Handler{
		paymentManager: paymentManager,
		auth:           auth,
		fileStorage:    fileStorage,
	}
}

type newPaymentRequest struct {
	ID       gopay.ID   `param:"id"`
	Template string     `json:"template"`
	User     gopay.User `json:"user"`
}

func (h Handler) NewPayment(c echo.Context) error {
	// admin auth
	log := slog.Default().With(
		slog.String("op", "Handler.NewPayment"),
		slog.String("request_id", c.Response().Header().Get(echo.HeaderXRequestID)),
	)

	var req newPaymentRequest
	if err := c.Bind(&req); err != nil {
		log.Error(err.Error())

		return c.String(http.StatusBadRequest, "invalid request")
	}

	// validate request

	// pass request id
	link, err := h.paymentManager.CreatePayment(req.Template, req.User)
	if err != nil {
		log.Error(err.Error())

		return c.String(http.StatusInternalServerError, "create payment failed")
	}

	if !link.IsValid() {
		log.Error("created link is invalid")

		return c.String(http.StatusInternalServerError, "created link is invalid")
	}

	log.Info("success payment created")

	return c.String(http.StatusOK, string(link))
}

func (h Handler) Redirect(c echo.Context) error {
	log := slog.Default().With(
		slog.String("op", "Handler.NewPayment"),
		slog.String("request_id", c.Response().Header().Get(echo.HeaderXRequestID)),
	)

	id := gopay.ID(c.Param("id"))
	if !id.IsValid() {
		log.Error("invalid request: bad id")

		return c.String(http.StatusBadRequest, "invalid request: bad id")
	}

	link, err := h.paymentManager.GetRedirectLink(id)
	if err != nil {
		log.Error(err.Error())

		return c.String(http.StatusInternalServerError, "get redirect link failed")
	}

	if !link.IsValid() {
		log.Error("redirect link is invalid")

		return c.String(http.StatusInternalServerError, "redirect link is invalid")
	}

	log.Info("success get redirect link")

	return c.Redirect(http.StatusTemporaryRedirect, string(link))
}

type checkoutRequest struct {
	ID     gopay.ID     `json:"id"`
	Status gopay.Status `json:"status"`
}

func (h Handler) Checkout(c echo.Context) error {
	// payment service auth
	log := slog.Default().With(
		slog.String("op", "Handler.NewPayment"),
		slog.String("request_id", c.Response().Header().Get(echo.HeaderXRequestID)),
	)

	var req checkoutRequest
	if err := c.Bind(&req); err != nil {
		log.Error(err.Error())

		return c.String(http.StatusBadRequest, "invalid request")
	}

	// validate request

	if err := h.paymentManager.UpdatePaymentStatus(req.ID, req.Status); err != nil {
		log.Error(err.Error())

		return c.String(http.StatusInternalServerError, "update payment status failed")
	}

	log.Info("success payment updated")

	return c.NoContent(http.StatusOK)
}

func (h Handler) File(c echo.Context) error {
	log := slog.Default().With(
		slog.String("op", "Handler.NewPayment"),
		slog.String("request_id", c.Response().Header().Get(echo.HeaderXRequestID)),
	)

	id := gopay.ID(c.Param("id"))
	if !id.IsValid() {
		log.Error("invalid request: bad id")

		return c.String(http.StatusBadRequest, "invalid request: bad id")
	}

	data, err := h.fileStorage.GetData(id)
	if err != nil {
		log.Error(err.Error())

		return c.String(http.StatusInternalServerError, "get file failed")
	}

	log.Info("success get file data")

	return c.Blob(http.StatusOK, "application/pdf", data)
}
