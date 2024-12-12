package bolt

import (
	"fmt"

	bolt "go.etcd.io/bbolt"

	"github.com/Anton-Kraev/gopay/models"
)

func (r PaymentRepository) SetLink(id models.ID, link models.Link) error {
	const op = "PaymentRepository.SetLink"

	if err := r.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(linkBucket)

		return b.Put([]byte(id), []byte(link))
	}); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r PaymentRepository) GetLink(id models.ID) (models.Link, error) {
	const op = "PaymentRepository.GetLink"

	var link models.Link

	if err := r.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(linkBucket)
		link = models.Link(b.Get([]byte(id)))

		return nil
	}); err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	if link == "" {
		return "", fmt.Errorf("%s: %w", op, errLinkNotFound)
	}

	return link, nil
}
