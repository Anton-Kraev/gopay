//go:generate mockgen -package=mocks -source=./gopay.go -destination=./mocks/gopay_mocks.go

package gopay

import "errors"

var ErrCreatePayment = errors.New("create payment failed")

type (
	templates interface {
		GetTemplate(templateName string) (PaymentTemplate, error)
	}

	linkGenerator interface {
		GenerateLink(id ID) (Link, error)
	}

	storage interface {
		SetUser(id ID, user User) error
		GetLinks(id ID) (Links, error)
		UpdateStatus(id ID, status Status) error
	}

	paymentService interface {
		CreatePayment(userID ID, template PaymentTemplate) (*Payment, error)
	}
)

type PaymentManager struct {
	templates templates
	links     linkGenerator
	storage   storage
	payments  paymentService
}

func NewPaymentManager(templates templates, linkGenerator linkGenerator, storage storage, paymentService paymentService) *PaymentManager {
	return &PaymentManager{
		templates: templates,
		links:     linkGenerator,
		storage:   storage,
		payments:  paymentService,
	}
}

func (pm *PaymentManager) NewPayment(templateName string, user User) (Link, error) {
	template, err := pm.templates.GetTemplate(templateName)
	if err != nil {
		return "", err
	}

	payment, err := pm.payments.CreatePayment(user.ID, template)
	if err != nil {
		return "", err
	}

	if payment == nil {
		return "", ErrCreatePayment
	}

	if err = pm.storage.SetUser(user.ID, user); err != nil {
		return "", err
	}

	return pm.links.GenerateLink(user.ID)
}

func (pm *PaymentManager) FollowLink(id ID) (Status, Link, error) {
	payment, err := pm.storage.GetLinks(id)
	if err != nil {
		return "", "", err
	}

	if payment.Status == StatusSucceeded {
		return payment.Status, payment.ResourceLink, nil
	}

	return payment.Status, payment.PaymentLink, nil
}

func (pm *PaymentManager) Checkout(id ID, newStatus Status) error {
	return pm.storage.UpdateStatus(id, newStatus)
}
