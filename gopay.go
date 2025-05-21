package gopay

import "errors"

var ErrCreatePayment = errors.New("create payment failed")

type (
	linkGenerator interface {
		GenerateLink() (ID, Link, error)
	}

	paymentStorage interface {
		Get(id ID) (Payment, error)
		Set(id ID, pay Payment) error
		GetStatus(id ID) (Status, error)
		GetStatuses() (map[ID]Status, error)
		UpdateStatus(id ID, status Status) error
		SetLink(id ID, link Link) error
		GetLink(id ID) (Link, error)
	}

	paymentService interface {
		CreatePayment(id ID, template PaymentTemplate) (*Payment, error)
	}
)

type PaymentManager struct {
	links    linkGenerator
	storage  paymentStorage
	payments paymentService
}

func NewPaymentManager(
	linkGenerator linkGenerator, paymentStorage paymentStorage, paymentService paymentService,
) *PaymentManager {
	return &PaymentManager{
		links:    linkGenerator,
		storage:  paymentStorage,
		payments: paymentService,
	}
}

func (pm *PaymentManager) CreatePayment(template PaymentTemplate, user User) (Link, error) {
	id, link, err := pm.links.GenerateLink()
	if err != nil {
		return "", err
	}

	payment, err := pm.payments.CreatePayment(id, template)
	if err != nil {
		return "", err
	}

	if payment == nil {
		return "", ErrCreatePayment
	}

	payment.User = user
	payment.ResourceLink = template.ResourceLink

	if err = pm.storage.Set(id, *payment); err != nil {
		return "", err
	}

	if err = pm.storage.SetLink(id, payment.PaymentLink); err != nil {
		return "", err
	}

	return link, nil
}

func (pm *PaymentManager) GetAllPaymentsStatuses() (map[ID]Status, error) {
	return pm.storage.GetStatuses()
}

func (pm *PaymentManager) GetPaymentStatus(id ID) (Status, error) {
	return pm.storage.GetStatus(id)
}

func (pm *PaymentManager) GetRedirectLink(id ID) (Link, error) {
	return pm.storage.GetLink(id)
}

func (pm *PaymentManager) UpdatePaymentStatus(id ID, newStatus Status) error {
	if newStatus == StatusSucceeded {
		payment, err := pm.storage.Get(id)
		if err != nil {
			return err
		}

		err = pm.storage.SetLink(id, payment.ResourceLink)
		if err != nil {
			return err
		}
	}

	return pm.storage.UpdateStatus(id, newStatus)
}
