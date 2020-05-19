package models

import (
	"time"
)

type Argue struct {
	ID           uint64    `json:"id" pg:",pk"`
	DealID       uint64    `json:"deal_id" pg:",notnull"`
	FirstUser    uint64    `json:"first_user"`
	SecondUser   uint64    `json:"second_user"`
	Category     uint8     `json:"category" pg:",notnull"`
	LastActivity time.Time `json:"last_activity"`
	Finished     bool      `json:"finished" pg:",use_zero"`
}
