package gopay

import (
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
)

type AdminClient interface {
	NewNewPaymentService() NewPaymentService
	NewAllPaymentService() AllPaymentService
	NewGetPaymentService() GetPaymentService
}

func NewAdminClient(serverURL string) AdminClient {
	return &adminClientImpl{api: resty.New().SetBaseURL(serverURL + "/api")}
}

type adminClientImpl struct {
	api *resty.Client
}

func (i *adminClientImpl) NewNewPaymentService() NewPaymentService {
	return &newPaymentServiceImpl{api: i.api}
}

func (i *adminClientImpl) NewAllPaymentService() AllPaymentService {
	return &allPaymentServiceImpl{api: i.api}
}

func (i *adminClientImpl) NewGetPaymentService() GetPaymentService {
	return &getPaymentServiceImpl{api: i.api}
}

type NewPaymentService interface {
	Currency(currency string) NewPaymentService
	Amount(amount uint) NewPaymentService
	Description(description string) NewPaymentService
	Do() error
}

type newPaymentServiceImpl struct {
	api         *resty.Client
	currency    string
	amount      uint
	description string
}

func (i *newPaymentServiceImpl) Currency(currency string) NewPaymentService {
	i.currency = currency

	return i
}

func (i *newPaymentServiceImpl) Amount(amount uint) NewPaymentService {
	i.amount = amount

	return i
}

func (i *newPaymentServiceImpl) Description(description string) NewPaymentService {
	i.description = description

	return i
}

func (i *newPaymentServiceImpl) Do() error {
	req := struct {
		currency    string
		amount      uint
		description string
	}{
		currency:    i.currency,
		amount:      i.amount,
		description: i.description,
	}

	// TODO: generate id more correctly and fix request body
	resp, err := i.api.R().SetBody(req).Post("/payments/123")
	if err != nil {
		return fmt.Errorf("AdminClient.NewPayment: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("AdminClient.NewPayment: error response from API %s", resp.String())
	}

	return nil
}

type AllPaymentService interface {
	Do() (map[ID]Status, error)
}

type allPaymentServiceImpl struct {
	api *resty.Client
}

func (i *allPaymentServiceImpl) Do() (map[ID]Status, error) {
	var res struct {
		Statuses []struct {
			ID     ID     `json:"id"`
			Status Status `json:"status"`
		} `json:"statuses"`
	}

	resp, err := i.api.R().SetResult(&res).Get("/payments")
	if err != nil {
		return nil, fmt.Errorf("AdminClient.AllPayment: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("AdminClient.AllPayment: error response from API %s", resp.String())
	}

	statuses := make(map[ID]Status)

	for _, payment := range res.Statuses {
		statuses[payment.ID] = payment.Status
	}

	return statuses, nil
}

type GetPaymentService interface {
	ID(id ID) GetPaymentService
	Do() (Status, error)
}

type getPaymentServiceImpl struct {
	api *resty.Client
	id  ID
}

func (i *getPaymentServiceImpl) ID(id ID) GetPaymentService {
	i.id = id

	return i
}

func (i *getPaymentServiceImpl) Do() (Status, error) {
	var status Status

	resp, err := i.api.R().SetResult(&status).Get("/payments/" + string(i.id))
	if err != nil {
		return "", fmt.Errorf("AdminClient.GetPayment: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return "", fmt.Errorf("AdminClient.GetPayment: error response from API %s", resp.String())
	}

	return status, nil
}
