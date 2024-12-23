package yookassa

type Amount struct {
	Value    string `json:"value"`
	Currency string `json:"currency"`
}

type Confirmation struct {
	Type      string `json:"type"`
	ReturnURL string `json:"return_url"`
}

type Metadata struct {
	OrderID string `json:"order_id"`
	UserID  string `json:"user_id"`
}

type Payment struct {
	ID           string       `json:"id"`
	Amount       Amount       `json:"amount"`
	Confirmation Confirmation `json:"confirmation"`
	Metadata     Metadata     `json:"metadata"`
	Description  string       `json:"description"`
	Capture      bool         `json:"capture"`
}
