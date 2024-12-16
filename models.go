package gopay

type ID string

type Status string

type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Link string

type Payment struct {
	Price  uint   `json:"price"`
	Status Status `json:"status"`
	User   User   `json:"user"`
}
