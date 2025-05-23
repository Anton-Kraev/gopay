package gopay

import (
	"net/url"
	"slices"

	"github.com/google/uuid"
)

type ID string

func (id ID) Validate() bool {
	return uuid.Validate(string(id)) == nil
}

type Status string

const (
	StatusPending           Status = "pending"
	StatusWaitingForCapture Status = "waiting_for_capture"
	StatusSucceeded         Status = "succeeded"
	StatusCancelled         Status = "canceled"
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
	ID    ID     `json:"id" validate:"required"`
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
	Currency     string `json:"currency" validate:"required"`
	Amount       uint   `json:"amount" validate:"required"`
	Description  string `json:"description" validate:"required"`
	ResourceLink Link   `json:"resource_link" validate:"required,url"`
}
