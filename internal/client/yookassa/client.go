package yookassa

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
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

func (c Client) CreatePayment(userID, description string, amount int) (*Payment, error) {
	const op = "Client.CreatePayment"

	uid := uuid.New().String()
	payment := &Payment{
		Amount: Amount{
			Value:    strconv.Itoa(amount),
			Currency: "RUB",
		},
		Confirmation: Confirmation{
			Type:      "redirect",
			ReturnURL: c.checkoutURL,
		},
		Metadata: Metadata{
			OrderID: uid,
			UserID:  userID,
		},
		Description: description,
		Capture:     true,
	}

	resp, err := c.http.R().
		SetBody(payment).
		SetResult(payment).
		SetHeader("Idempotence-Key", uid).
		Post("/payment")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("%s: error response from API %s", op, resp.String())
	}

	if payment.ID == "" {
		return nil, fmt.Errorf("%s: empty payment ID", op)
	}

	return payment, nil
}
