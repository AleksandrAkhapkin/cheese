package types

type PayRes struct {
	Request_id      string  `json:"request_id"`
	Contract_amount float64 `json:"contract_amount"`
	Title           string  `json:"title"`
	Error           string  `json:"error"`
	Status          string  `json:"status"`
}
type PayedRes struct {
	Status     string `json:"status"`
	Payment_id string `json:"payment_id"`
	Invoice_id string `json:"invoice_id"`
	Error      string `json:"error"`
}
