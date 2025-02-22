package handler

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/Anton-Kraev/gopay"
)

type fileStorage interface {
	GetData(id gopay.ID) ([]byte, error)
}

type Handler struct {
	paymentManager *gopay.PaymentManager
	fileStorage    fileStorage
}

func NewHandler(paymentManager *gopay.PaymentManager, fileStorage fileStorage) Handler {
	return Handler{
		paymentManager: paymentManager,
		fileStorage:    fileStorage,
	}
}

type newPaymentRequest struct {
	ID       gopay.ID   `param:"id" validate:"required,id"`
	Template string     `json:"template" validate:"required"`
	User     gopay.User `json:"user"`
}

func (h Handler) NewPayment(c echo.Context) error {
	log := slog.Default().With(
		slog.String("op", "Handler.NewPayment"),
		slog.String("request_id", c.Response().Header().Get(echo.HeaderXRequestID)),
	)

	var req newPaymentRequest
	if err := c.Bind(&req); err != nil {
		log.Error(err.Error())

		return c.String(http.StatusBadRequest, "invalid request")
	}

	if err := c.Validate(&req); err != nil {
		log.Error(err.Error())

		return c.String(http.StatusBadRequest, "invalid request")
	}

	link, err := h.paymentManager.CreatePayment(req.ID, req.Template, req.User)
	if err != nil {
		log.Error(err.Error())

		return c.String(http.StatusInternalServerError, "create payment failed")
	}

	if !link.Validate() {
		log.Error("created link is invalid")

		return c.String(http.StatusInternalServerError, "created link is invalid")
	}

	log.Info("success payment created")

	return c.String(http.StatusOK, string(link))
}

func (h Handler) Redirect(c echo.Context) error {
	log := slog.Default().With(
		slog.String("op", "Handler.Redirect"),
		slog.String("request_id", c.Response().Header().Get(echo.HeaderXRequestID)),
	)

	id := gopay.ID(c.Param("id"))
	if !id.Validate() {
		log.Error("invalid request: bad id")

		return c.String(http.StatusBadRequest, "invalid request: bad id")
	}

	link, err := h.paymentManager.GetRedirectLink(id)
	if err != nil {
		log.Error(err.Error())

		return c.String(http.StatusInternalServerError, "get redirect link failed")
	}

	if !link.Validate() {
		log.Error("redirect link is invalid")

		return c.String(http.StatusInternalServerError, "redirect link is invalid")
	}

	log.Info("success get redirect link")

	return c.Redirect(http.StatusTemporaryRedirect, string(link))
}

type checkoutRequest struct {
	ID     gopay.ID `param:"id" validate:"required,id"`
	Object struct {
		Status gopay.Status `json:"status" validate:"required,status"`
	} `json:"object"`
}

func (h Handler) Checkout(c echo.Context) error {
	log := slog.Default().With(
		slog.String("op", "Handler.Checkout"),
		slog.String("request_id", c.Response().Header().Get(echo.HeaderXRequestID)),
	)

	var req checkoutRequest
	if err := c.Bind(&req); err != nil {
		log.Error(err.Error())

		return c.String(http.StatusBadRequest, "invalid request")
	}

	if err := c.Validate(&req); err != nil {
		log.Error(err.Error())

		return c.String(http.StatusBadRequest, "invalid request")
	}

	if err := h.paymentManager.UpdatePaymentStatus(req.ID, req.Object.Status); err != nil {
		log.Error(err.Error())

		return c.String(http.StatusInternalServerError, "update payment status failed")
	}

	log.Info("success payment updated")

	return c.NoContent(http.StatusOK)
}

func (h Handler) File(c echo.Context) error {
	log := slog.Default().With(
		slog.String("op", "Handler.File"),
		slog.String("request_id", c.Response().Header().Get(echo.HeaderXRequestID)),
	)

	id := gopay.ID(c.Param("id"))
	if !id.Validate() {
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
