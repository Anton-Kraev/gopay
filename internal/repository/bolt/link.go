package bolt

import (
	"context"

	"github.com/Anton-Kraev/gopay/internal/models/payment"
)

func (r PaymentRepository) SetLink(ctx context.Context, id payment.ID, link payment.Link) error {
	panic("not implemented")
}

func (r PaymentRepository) GetLink(ctx context.Context, id payment.ID) (payment.Link, error) {
	panic("not implemented")
}
