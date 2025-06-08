package gopay

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/go-resty/resty/v2"
)

type AdminClient interface {
	NewNewPaymentService() NewPaymentService
	NewAllPaymentService() AllPaymentService
	NewGetPaymentService() GetPaymentService
}

func NewAdminClient(serverURL string) (AdminClient, error) {
	baseURL, err := url.ParseRequestURI(serverURL)
	if err != nil {
		return nil, fmt.Errorf("gopay.NewAdminClient: %w", err)
	}

	apiURL := baseURL.JoinPath("api").String()

	return &adminClientImpl{api: resty.New().SetBaseURL(apiURL)}, nil
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
	ResourceLink(link Link) NewPaymentService
	Do() (Link, error)

	String() string
}

type newPaymentServiceImpl struct {
	api         *resty.Client
	currency    string
	amount      uint
	description string
	link        Link
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

func (i *newPaymentServiceImpl) ResourceLink(link Link) NewPaymentService {
	i.link = link

	return i
}

type newPaymentRequest struct {
	Template PaymentTemplate `json:"template"`
	User     User            `json:"user"`
}

func (i *newPaymentServiceImpl) Do() (Link, error) {
	if !i.link.Validate() {
		return "", fmt.Errorf("AdminClient.NewPayment: invalid link %s", i.link)
	}

	req := newPaymentRequest{
		Template: PaymentTemplate{
			Currency:     i.currency,
			Amount:       i.amount,
			Description:  i.description,
			ResourceLink: i.link,
		},
		User: User{
			ID:    "id",
			Name:  "name",
			Email: "email@mail.com",
		},
	}

	resp, err := i.api.R().SetBody(&req).Post("/payments")
	if err != nil {
		return "", fmt.Errorf("AdminClient.NewPayment: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return "", fmt.Errorf("AdminClient.NewPayment: error response from API %s", resp.String())
	}

	return Link(resp.String()), nil
}

func (i *newPaymentServiceImpl) String() string {
	return fmt.Sprintf(
		"сумма: %d\nвалюта: %s\nописание: %s\nссылка на ресурс: %s",
		i.amount, i.currency, i.description, i.link,
	)
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
	if !i.id.Validate() {
		return "", fmt.Errorf("AdminClient.GetPayment: invalid id %s", i.id)
	}

	resp, err := i.api.R().Get(fmt.Sprintf("/payments/%s", i.id))
	if err != nil {
		return "", fmt.Errorf("AdminClient.GetPayment: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return "", fmt.Errorf("AdminClient.GetPayment: error response from API %s", resp.String())
	}

	return Status(resp.String()), nil
}
