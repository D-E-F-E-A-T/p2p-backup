package models

import (
	"time"
)

type Chat struct {
	ID           uint64 `json:"id" pg:",pk"`
	DealID       uint64 `json:"deal_id" pg:",pk"`
	FirstUser    uint64 `json:"first_user"`
	SecondUser   uint64 `json:"second_user"`
	FirstUnread  uint   `json:"first_unread"`
	SecondUnread uint   `json:"second_unread"`
	ArgueID      uint64 `json:"argue_id"`
	Closed       bool   `json:"closed" pg:",use_zero"`
	// TODO update timestamp on every message
	Timestamp time.Time `json:"timestamp"`
}
