package bolt

import (
	"encoding/json"
	"fmt"

	bolt "go.etcd.io/bbolt"

	"github.com/Anton-Kraev/gopay/models"
)

func (r PaymentRepository) Get(id models.ID) (models.Payment, error) {
	const op = "PaymentRepository.Get"

	var pay models.Payment

	if err := r.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(paymentBucket)

		binPay := b.Get([]byte(id))
		if len(binPay) == 0 {
			return errPaymentNotFound
		}

		return json.Unmarshal(binPay, &pay)
	}); err != nil {
		return models.Payment{}, fmt.Errorf("%s: %w", op, err)
	}

	return pay, nil
}

func (r PaymentRepository) Set(id models.ID, pay models.Payment) error {
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

func (r PaymentRepository) UpdateStatus(id models.ID, status models.Status) error {
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
