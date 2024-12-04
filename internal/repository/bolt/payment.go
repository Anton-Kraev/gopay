package bolt

import (
	"encoding/json"
	"fmt"

	bolt "go.etcd.io/bbolt"

	"github.com/Anton-Kraev/gopay/internal/models/payment"
)

func (r PaymentRepository) Get(id payment.ID) (payment.Payment, error) {
	const op = "PaymentRepository.Get"

	var pay payment.Payment

	if err := r.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(paymentBucket)

		binPay := b.Get([]byte(id))
		if len(binPay) == 0 {
			return errPaymentNotFound
		}

		return json.Unmarshal(binPay, &pay)
	}); err != nil {
		return payment.Payment{}, fmt.Errorf("%s: %w", op, err)
	}

	return pay, nil
}

func (r PaymentRepository) Set(id payment.ID, pay payment.Payment) error {
	const op = "PaymentRepository.SetLink"

	if err := r.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(paymentBucket)

		binPay, err := json.Marshal(pay)
		if err != nil {
			return err
		}

		return b.Put([]byte(id), binPay)
	}); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r PaymentRepository) UpdateStatus(id payment.ID, status payment.Status) error {
	const op = "PaymentRepository.UpdateStatus"

	pay, err := r.Get(id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	pay.Status = status

	if err = r.Set(id, pay); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
