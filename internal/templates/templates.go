package templates

import (
	"encoding/json"
	"errors"
	"fmt"

	bolt "go.etcd.io/bbolt"

	"github.com/Anton-Kraev/gopay"
)

var (
	templateBucket      = []byte("TemplateBucket")
	errTemplateNotFound = errors.New("template not found")
)

type Templates struct {
	db *bolt.DB
}

func New(db *bolt.DB) (Templates, error) {
	if err := createBucket(db); err != nil {
		return Templates{}, fmt.Errorf("templates.New: %w", err)
	}

	return Templates{db: db}, nil
}

func createBucket(db *bolt.DB) error {
	if err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(templateBucket)

		return err
	}); err != nil {
		return fmt.Errorf("templates.createBucket: %w", err)
	}

	return nil
}

func (t Templates) GetTemplate(templateName string) (gopay.PaymentTemplate, error) {
	var paymentTemplate gopay.PaymentTemplate

	if err := t.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(templateBucket)

		binPay := b.Get([]byte(templateName))
		if len(binPay) == 0 {
			return errTemplateNotFound
		}

		return json.Unmarshal(binPay, &paymentTemplate)
	}); err != nil {
		return gopay.PaymentTemplate{}, fmt.Errorf("templates.Templates.GetTemplate: %w", err)
	}

	return paymentTemplate, nil
}

func (t Templates) SetTemplate(templateName string, template gopay.PaymentTemplate) error {
	if err := t.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(templateBucket)

		binTemplate, err := json.Marshal(template)
		if err != nil {
			return err
		}

		return b.Put([]byte(templateName), binTemplate)
	}); err != nil {
		return fmt.Errorf("templates.Templates.SetTemplate: %w", err)
	}

	return nil
}
