package models

type Wallet struct {
	UserID   uint64  `json:"user_id" pg:",nopk,notnull,unique:user_currency"`
	Currency uint8   `json:"currency" pg:",notnull,unique:user_currency"`
	Balance  float32 `json:"balance"`
}
