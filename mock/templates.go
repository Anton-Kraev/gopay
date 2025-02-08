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
		PaymentLink:  "https://github.com/Anton-Kraev/gopay",
		ResourceLink: "http://127.0.0.1:1323/file/123",
	}, nil
}
