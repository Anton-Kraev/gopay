package bolt

import (
	"fmt"

	bolt "go.etcd.io/bbolt"

	"github.com/Anton-Kraev/gopay"
)

func (r PaymentRepository) SetLink(id gopay.ID, link gopay.Link) error {
	if err := r.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(linkBucket)

		return b.Put([]byte(id), []byte(link))
	}); err != nil {
		return fmt.Errorf("bolt.PaymentRepository.SetLink: %w", err)
	}

	return nil
}

func (r PaymentRepository) GetLink(id gopay.ID) (gopay.Link, error) {
	var link gopay.Link

	if err := r.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(linkBucket)

		link = gopay.Link(b.Get([]byte(id)))
		if link == "" {
			return errLinkNotFound
		}

		return nil
	}); err != nil {
		return "", fmt.Errorf("bolt.PaymentRepository.GetLink: %w", err)
	}

	return link, nil
}
