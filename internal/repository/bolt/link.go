package bolt

import (
	"fmt"

	bolt "go.etcd.io/bbolt"

	"github.com/Anton-Kraev/gopay/internal/models/payment"
)

func (r PaymentRepository) SetLink(id payment.ID, link payment.Link) error {
	const op = "PaymentRepository.SetLink"

	if err := r.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(linkBucket)

		return b.Put([]byte(id), []byte(link))
	}); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r PaymentRepository) GetLink(id payment.ID) (payment.Link, error) {
	const op = "PaymentRepository.GetLink"

	var link payment.Link

	if err := r.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(linkBucket)
		link = payment.Link(b.Get([]byte(id)))

		return nil
	}); err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	if link == "" {
		return "", fmt.Errorf("%s: %w", op, errLinkNotFound)
	}

	return link, nil
}
