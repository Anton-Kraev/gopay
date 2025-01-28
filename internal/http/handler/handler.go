package handler

import (
	"github.com/labstack/echo/v4"

	"github.com/Anton-Kraev/gopay"
)

type (
	auth        interface{}
	fileStorage interface{}
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

func (h Handler) NewPayment(c echo.Context) error {
	// logs and metrics

	// admin auth
	// get id from path
	// get user and template from body
	// validate request
	// create payment with payment manager
	// return link

	panic("implement me")
}

func (h Handler) Redirect(c echo.Context) error {
	// logs and metrics

	// get id from path
	// get link from payment manager
	// redirect to link

	panic("implement me")
}

func (h Handler) Checkout(c echo.Context) error {
	// logs and metrics

	// payment service auth
	// get id from path
	// get new status from body
	// validate request
	// update status with payment manager

	panic("implement me")
}

func (h Handler) File(c echo.Context) error {
	// logs and metrics

	// get id from path
	// get file from external file storage
	// return file with appropriate headers

	panic("implement me")
}
