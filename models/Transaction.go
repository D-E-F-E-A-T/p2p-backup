package models

import (
	"time"
)

type Transaction struct {
	ID        uint64    `json:"id" pg:",nopk"`
	UserID    uint64    `json:"user"`
	Deal      uint64    `json:"deal"`
	Currency  uint8     `json:"currency" pg:",notnull"`
	Amount    float32   `json:"amount" pg:",notnull"`
	Timestamp time.Time `json:"timestamp" pg:",notnull"`
}
