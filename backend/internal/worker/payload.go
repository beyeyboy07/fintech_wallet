package worker

type TopUpPayload struct {
	UserID        uint  `json:"user_id"`
	TransactionID uint  `json:"transaction_id"`
	Amount        int64 `json:"amount"`
}
