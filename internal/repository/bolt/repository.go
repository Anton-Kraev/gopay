package bolt

import bolt "go.etcd.io/bbolt"

var (
	paymentBucket = []byte("PaymentBucket")
	linkBucket    = []byte("LinkBucket")
)

type PaymentRepository struct {
	db *bolt.DB
}

func NewPaymentRepository(db *bolt.DB) PaymentRepository {
	return PaymentRepository{db: db}
}
