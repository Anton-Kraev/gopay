package bolt

import (
	"encoding/json"
	"fmt"

	bolt "go.etcd.io/bbolt"

	"github.com/Anton-Kraev/gopay"
)

func (r PaymentRepository) Get(id gopay.ID) (gopay.Payment, error) {
	var pay gopay.Payment

	if err := r.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(paymentBucket)

		binPay := b.Get([]byte(id))
		if len(binPay) == 0 {
			return errPaymentNotFound
		}

		return json.Unmarshal(binPay, &pay)
	}); err != nil {
		return gopay.Payment{}, fmt.Errorf("bolt.PaymentRepository.Get: %w", err)
	}

	return pay, nil
}

func (r PaymentRepository) Set(id gopay.ID, pay gopay.Payment) error {
	if err := r.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(paymentBucket)

		binPay, err := json.Marshal(pay)
		if err != nil {
			return err
		}

		return b.Put([]byte(id), binPay)
	}); err != nil {
		return fmt.Errorf("bolt.PaymentRepository.Set: %w", err)
	}

	return nil
}

func (r PaymentRepository) UpdateStatus(id gopay.ID, status gopay.Status) error {
	const op = "bolt.PaymentRepository.UpdateStatus"

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
