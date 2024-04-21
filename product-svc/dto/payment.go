package dto

type Payment struct {
	ID           uint64 `json:"id"`
	UserID       uint64 `json:"user_id"`
	CurrencyCode string `json:"currency_code"`
	Amount       int64  `json:"amount"`
}
