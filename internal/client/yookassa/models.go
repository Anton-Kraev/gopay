package yookassa

type Amount struct {
	Value    string `json:"value"`
	Currency string `json:"currency"`
}

type Confirmation struct {
	Type            string `json:"type"`
	ReturnURL       string `json:"return_url"`
	ConfirmationURL string `json:"confirmation_url"`
}

type Metadata struct {
	ID string `json:"id"`
}

type Payment struct {
	ID           string       `json:"id"`
	Status       string       `json:"status"`
	Amount       Amount       `json:"amount"`
	Confirmation Confirmation `json:"confirmation"`
	Metadata     Metadata     `json:"metadata"`
	Description  string       `json:"description"`
	Capture      bool         `json:"capture"`
}
