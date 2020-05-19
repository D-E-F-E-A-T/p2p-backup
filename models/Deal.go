package models

import (
	"time"
)

type Deal struct {
	ID           uint64    `json:"id" pg:",pk"`
	OfferID      uint64    `json:"offer_id" pg:",notnull"`
	UserID       uint64    `json:"user_id" pg:",notnull"`
	Amount       float32   `json:"amount" pg:",notnull"`
	ServiceFee   float32   `json:"service_fee"`
	FixedAmount  float32   `json:"fixed_amount"`
	Timestamp    time.Time `json:"timestamp"`
	Accepted     bool      `json:"accepted" pg:",use_zero"`
	FromApproved bool      `json:"from_approved" pg:",use_zero"`
	ToApproved   bool      `json:"to_approved" pg:",use_zero"`
	Finished     bool      `json:"finished" pg:",use_zero"`
}
