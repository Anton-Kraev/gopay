package bolt

import (
	"context"

	"github.com/Anton-Kraev/gopay/internal/models/payment"
)

func (r PaymentRepository) Get(ctx context.Context, id payment.ID) (payment.Payment, error) {
	panic("not implemented")
}

func (r PaymentRepository) Set(ctx context.Context, id payment.ID, pay payment.Payment) error {
	panic("not implemented")
}

func (r PaymentRepository) UpdateStatus(ctx context.Context, status payment.Status) error {
	panic("not implemented")
}
