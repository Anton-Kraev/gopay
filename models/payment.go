package models

type ID string

type Status string

type Payment struct {
	Price  uint   `json:"price"`
	Status Status `json:"status"`
	User   User   `json:"user"`
}

type Link string
