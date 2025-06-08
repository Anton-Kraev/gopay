package telegram

type state string

const (
	stateNewPaymentAmount       state = "new_payment_amount"
	stateNewPaymentDescription  state = "new_payment_description"
	stateNewPaymentLink         state = "new_payment_link"
	stateNewPaymentConfirmation state = "new_payment_confirmation"
)
