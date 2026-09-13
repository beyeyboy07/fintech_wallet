package model

type PaymentWebhookRequest struct {
	EventID       string `json:"event_id" binding:"required"`
	TransactionID uint   `json:"transaction_id" binding:"required"`
	Status        string `json:"status" binding:"required"`
	Amount        int64  `json:"amount" binding:"required"`
}
