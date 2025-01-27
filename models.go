package gopay

type ID string

type Status string

const (
	StatusPending           Status = "pending"
	StatusWaitingForCapture Status = "waiting_for_capture"
	StatusSucceeded         Status = "succeeded"
	StatusCancelled         Status = "cancelled"
)

type User struct {
	ID    ID     `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Link string

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
