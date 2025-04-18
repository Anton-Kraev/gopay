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
	Template gopay.PaymentTemplate `json:"template" validate:"required"`
	User     gopay.User            `json:"user" validate:"required"`
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

	link, err := h.paymentManager.CreatePayment(req.Template, req.User)
	if err != nil {
		log.Error(err.Error())

		return c.String(http.StatusInternalServerError, "create payment failed")
	}

	log.Info("success payment created")

	return c.String(http.StatusOK, string(link))
}

type paymentStatus struct {
	ID     gopay.ID     `json:"id"`
	Status gopay.Status `json:"status"`
}

type allPaymentResponse struct {
	Statuses []paymentStatus `json:"statuses"`
}

func (h Handler) AllPayment(c echo.Context) error {
	log := slog.Default().With(
		slog.String("op", "Handler.AllPayment"),
		slog.String("request_id", c.Response().Header().Get(echo.HeaderXRequestID)),
	)

	statuses, err := h.paymentManager.GetAllPaymentsStatuses()
	if err != nil {
		log.Error(err.Error())

		return c.String(http.StatusInternalServerError, "get payments statuses failed")
	}

	log.Info("success get payments statuses")

	var resp allPaymentResponse

	for id, status := range statuses {
		resp.Statuses = append(resp.Statuses, paymentStatus{
			ID:     id,
			Status: status,
		})
	}

	return c.JSON(http.StatusOK, resp)
}

func (h Handler) GetPayment(c echo.Context) error {
	log := slog.Default().With(
		slog.String("op", "Handler.GetPayment"),
		slog.String("request_id", c.Response().Header().Get(echo.HeaderXRequestID)),
	)

	id := gopay.ID(c.Param("id"))
	if !id.Validate() {
		log.Error("invalid request: bad id")

		return c.String(http.StatusBadRequest, "invalid request: bad id")
	}

	status, err := h.paymentManager.GetPaymentStatus(id)
	if err != nil {
		log.Error(err.Error())

		return c.String(http.StatusInternalServerError, "get payment status failed")
	}

	log.Info("success get payment status")

	return c.String(http.StatusOK, string(status))
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

	log.Info("success get redirect link")

	return c.Redirect(http.StatusTemporaryRedirect, string(link))
}

type checkoutRequest struct {
	Object struct {
		Metadata struct {
			ID gopay.ID `json:"id" validate:"required,id"`
		} `json:"metadata"`
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

	if err := h.paymentManager.UpdatePaymentStatus(req.Object.Metadata.ID, req.Object.Status); err != nil {
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

	data, err := h.fileStorage.GetData(id)
	if err != nil {
		log.Error(err.Error())

		return c.String(http.StatusInternalServerError, "get file failed")
	}

	log.Info("success get file data")

	return c.Blob(http.StatusOK, "application/pdf", data)
}
