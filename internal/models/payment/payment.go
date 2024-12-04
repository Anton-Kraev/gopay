package payment

import (
	"github.com/Anton-Kraev/gopay/internal/models/user"
)

type ID string

type Status string

type Payment struct {
	Price  uint      `json:"price"`
	Status Status    `json:"status"`
	User   user.User `json:"user"`
}

type Link string
