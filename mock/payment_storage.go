package mock

import (
	"errors"

	"github.com/Anton-Kraev/gopay"
)

var ErrNotFound = errors.New("not found")

type PaymentStorage struct {
	paymentStorage map[gopay.ID]gopay.Payment
	linkStorage    map[gopay.ID]gopay.Link
}

func NewPaymentStorage() *PaymentStorage {
	return &PaymentStorage{
		paymentStorage: make(map[gopay.ID]gopay.Payment),
		linkStorage:    make(map[gopay.ID]gopay.Link),
	}
}

func (p *PaymentStorage) Get(id gopay.ID) (gopay.Payment, error) {
	if payment, ok := p.paymentStorage[id]; ok {
		return payment, nil
	}

	return gopay.Payment{}, ErrNotFound
}

func (p *PaymentStorage) Set(id gopay.ID, pay gopay.Payment) error {
	p.paymentStorage[id] = pay

	return nil
}

func (p *PaymentStorage) UpdateStatus(id gopay.ID, status gopay.Status) error {
	payment, err := p.Get(id)
	if err != nil {
		return err
	}

	payment.Status = status

	return p.Set(id, payment)
}

func (p *PaymentStorage) SetLink(id gopay.ID, link gopay.Link) error {
	p.linkStorage[id] = link

	return nil
}

func (p *PaymentStorage) GetLink(id gopay.ID) (gopay.Link, error) {
	if link, ok := p.linkStorage[id]; ok {
		return link, nil
	}

	return "", ErrNotFound
}
