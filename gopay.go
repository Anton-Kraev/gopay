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

	paymentStorage interface {
		Get(id ID) (Payment, error)
		Set(id ID, pay Payment) error
		UpdateStatus(id ID, status Status) error
		SetLink(id ID, link Link) error
		GetLink(id ID) (Link, error)
	}

	paymentService interface {
		CreatePayment(userID ID, template PaymentTemplate) (*Payment, error)
	}
)

type PaymentManager struct {
	templates templates
	links     linkGenerator
	storage   paymentStorage
	payments  paymentService
}

func NewPaymentManager(
	templates templates, linkGenerator linkGenerator, paymentStorage paymentStorage, paymentService paymentService,
) *PaymentManager {
	return &PaymentManager{
		templates: templates,
		links:     linkGenerator,
		storage:   paymentStorage,
		payments:  paymentService,
	}
}

func (pm *PaymentManager) CreatePayment(id ID, templateName string, user User) (Link, error) {
	template, err := pm.templates.GetTemplate(templateName)
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
	payment.PaymentLink = template.PaymentLink
	payment.ResourceLink = template.ResourceLink

	if err = pm.storage.Set(id, *payment); err != nil {
		return "", err
	}

	if err = pm.storage.SetLink(id, payment.PaymentLink); err != nil {
		return "", err
	}

	return pm.links.GenerateLink(id)
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
