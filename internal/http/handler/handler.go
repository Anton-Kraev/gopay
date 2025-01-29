package handler

import (
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
	// logs and metrics and wrap error
	// admin auth

	var req newPaymentRequest
	if err := c.Bind(&req); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	// validate request

	// pass request id
	link, err := h.paymentManager.CreatePayment(req.Template, req.User)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	if !link.IsValid() {
		return c.String(http.StatusInternalServerError, "created link is invalid")
	}

	return c.String(http.StatusOK, string(link))
}

func (h Handler) Redirect(c echo.Context) error {
	// logs and metrics and wrap error

	id := gopay.ID(c.Param("id"))
	if !id.IsValid() {
		return c.String(http.StatusBadRequest, "id is invalid")
	}

	link, err := h.paymentManager.GetRedirectLink(id)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	if !link.IsValid() {
		return c.String(http.StatusInternalServerError, "redirect link is invalid")
	}

	return c.Redirect(http.StatusTemporaryRedirect, string(link))
}

type checkoutRequest struct {
	ID     gopay.ID     `json:"id"`
	Status gopay.Status `json:"status"`
}

func (h Handler) Checkout(c echo.Context) error {
	// logs and metrics and wrap error
	// payment service auth

	var req checkoutRequest
	if err := c.Bind(&req); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	// validate request

	if err := h.paymentManager.UpdatePaymentStatus(req.ID, req.Status); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusOK)
}

func (h Handler) File(c echo.Context) error {
	// logs and metrics and wrap error

	id := gopay.ID(c.Param("id"))
	if !id.IsValid() {
		return c.String(http.StatusBadRequest, "id is invalid")
	}

	data, err := h.fileStorage.GetData(id)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.Blob(http.StatusOK, "application/pdf", data)
}
