package yookassa

import (
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"

	"github.com/Anton-Kraev/gopay"
)

const (
	baseURL               = "https://api.yookassa.ru/v3"
	createPaymentEndpoint = "/payments"
)

type Client struct {
	checkoutURL string
	http        *resty.Client
}

func NewClient(checkoutURL string, config AuthConfig) Client {
	return Client{
		checkoutURL: checkoutURL,
		http: resty.New().
			SetBaseURL(baseURL).
			SetBasicAuth(config.ID, config.Token),
	}
}

func (c Client) CreatePayment(id gopay.ID, template gopay.PaymentTemplate) (*gopay.Payment, error) {
	const op = "yookassa.Client.CreatePayment"

	yookassaPayment := &Payment{
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
		SetBody(yookassaPayment).
		SetResult(yookassaPayment).
		SetHeader("Idempotence-Key", uuid.New().String()).
		Post(createPaymentEndpoint)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("%s: error response from API %s", op, resp.String())
	}

	payment := &gopay.Payment{
		Amount:      template.Amount,
		Status:      gopay.Status(yookassaPayment.Status),
		PaymentLink: gopay.Link(yookassaPayment.Confirmation.ConfirmationURL),
	}

	if yookassaPayment.ID == "" {
		return nil, fmt.Errorf("%s: empty payment ID", op)
	}

	if !payment.Status.Validate() {
		return nil, fmt.Errorf("%s: unknown payment status", op)
	}

	if !payment.PaymentLink.Validate() {
		return nil, fmt.Errorf("%s: bad payment url", op)
	}

	return payment, nil
}
