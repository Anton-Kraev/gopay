package gopay

import (
	"net/url"
	"slices"
)

type ID string

func (id ID) Validate() bool {
	return id != ""
}

type Status string

const (
	StatusPending           Status = "pending"
	StatusWaitingForCapture Status = "waiting_for_capture"
	StatusSucceeded         Status = "succeeded"
	StatusCancelled         Status = "cancelled"
)

func (s Status) Validate() bool {
	return slices.Contains([]Status{
		StatusPending,
		StatusWaitingForCapture,
		StatusSucceeded,
		StatusCancelled,
	}, s)
}

type User struct {
	ID    ID     `json:"id" validate:"required,id"`
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
}

type Link string

func (l Link) Validate() bool {
	_, err := url.ParseRequestURI(string(l))

	return err == nil
}

type Payment struct {
	User         User   `json:"user"`
	Amount       uint   `json:"amount"`
	Status       Status `json:"status"`
	PaymentLink  Link   `json:"payment_link"`
	ResourceLink Link   `json:"resource_link"`
}

type PaymentTemplate struct {
	Currency     string `json:"currency"`
	Amount       uint   `json:"amount"`
	Description  string `json:"description"`
	PaymentLink  Link   `json:"payment_link"`
	ResourceLink Link   `json:"resource_link"`
}
