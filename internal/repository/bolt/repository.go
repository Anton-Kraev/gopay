package bolt

import (
	"errors"
	"fmt"

	bolt "go.etcd.io/bbolt"
)

var (
	paymentBucket = []byte("PaymentBucket")
	linkBucket    = []byte("LinkBucket")

	errPaymentNotFound = errors.New("payment not found")
	errLinkNotFound    = errors.New("link not found")
)

type PaymentRepository struct {
	db *bolt.DB
}

func NewPaymentRepository(db *bolt.DB) (PaymentRepository, error) {
	const op = "NewPaymentRepository"

	if err := createBuckets(db); err != nil {
		return PaymentRepository{}, fmt.Errorf("%s: %w", op, err)
	}

	return PaymentRepository{db: db}, nil
}

func createBuckets(db *bolt.DB) error {
	return db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(paymentBucket)
		if err != nil {
			return err
		}

		_, err = tx.CreateBucketIfNotExists(linkBucket)

		return err
	})
}
