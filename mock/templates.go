package mock

import (
	"github.com/Anton-Kraev/gopay"
)

type Templates struct{}

func NewTemplates() Templates {
	return Templates{}
}

func (t Templates) GetTemplate(_ string) (gopay.PaymentTemplate, error) {
	return gopay.PaymentTemplate{
		Currency:     "RUB",
		Amount:       100,
		Description:  "description",
		ResourceLink: "http://127.0.0.1:8080/api/files/123",
	}, nil
}
