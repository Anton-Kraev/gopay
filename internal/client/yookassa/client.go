package yookassa

import (
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"

	"github.com/Anton-Kraev/gopay"
)

type Client struct {
	checkoutURL string
	http        *resty.Client
}

func NewClient(checkoutURL string, config Config) Client {
	return Client{
		checkoutURL: checkoutURL,
		http: resty.New().
			SetBaseURL(config.URL).
			SetBasicAuth(config.ID, config.Token),
	}
}

func (c Client) CreatePayment(id gopay.ID, template gopay.PaymentTemplate) (*gopay.Payment, error) {
	const op = "Client.CreatePayment"

	uid := uuid.New().String()
	payment := &Payment{
		Amount: Amount{
			Value:    fmt.Sprintf("%d", template.Amount),
			Currency: template.Currency,
		},
		Confirmation: Confirmation{
			Type:      "redirect",
			ReturnURL: c.checkoutURL,
		},
		Metadata: Metadata{
			ID: string(id),
		},
		Description: template.Description,
		Capture:     true,
	}

	resp, err := c.http.R().
		SetBody(payment).
		SetResult(payment).
		SetHeader("Idempotence-Key", uid).
		Post("/payments")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("%s: error response from API %s", op, resp.String())
	}

	if payment.ID == "" {
		return nil, fmt.Errorf("%s: empty payment ID", op)
	}

	return &gopay.Payment{
		Amount:      template.Amount,
		Status:      gopay.Status(payment.Status),
		PaymentLink: gopay.Link(payment.Confirmation.ConfirmationURL),
	}, nil
}
