package mock

import (
	"github.com/Anton-Kraev/gopay"
)

type PaymentService struct{}

func NewPaymentService() PaymentService {
	return PaymentService{}
}

func (p PaymentService) CreatePayment(_ gopay.ID, template gopay.PaymentTemplate) (*gopay.Payment, error) {
	return &gopay.Payment{
		Amount: template.Amount,
		Status: gopay.StatusWaitingForCapture,
	}, nil
}
