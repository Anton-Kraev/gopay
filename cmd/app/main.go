package main

import (
	"log"
	"time"

	bolt "go.etcd.io/bbolt"

	repo "github.com/Anton-Kraev/gopay/internal/repository/bolt"
)

const (
	dbFilePath    = "my.db"
	dbOpenTimeout = 1 * time.Second
)

func main() {
	db, err := bolt.Open(
		dbFilePath,
		0600,
		&bolt.Options{Timeout: dbOpenTimeout},
	)
	if err != nil {
		log.Fatalln(err)
	}

	defer func(db *bolt.DB) {
		if err = db.Close(); err != nil {
			log.Println("error closing db connection:", err)

			return
		}

		log.Println("db connection closed")
	}(db)

	storage := repo.NewPaymentRepository(db)
	_ = storage

	log.Println("db connection opened")
}
