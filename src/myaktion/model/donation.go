package model

type Donation struct {
	Amount float64 `json:"amount"`
	ReceiptRequested bool `json:"receiptRequested"`
	DonorName string `json:"donorName"`
	Status Status `json:"status"`
	Account Account `json:"account"`
}

type Status string

const (
	TRANSFERRED Status = "TRANSFERRED"
	IN_PROCESS Status = "IN_PROCESS"
)

