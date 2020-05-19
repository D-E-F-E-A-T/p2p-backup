package models

import (
	"time"
)

type Message struct {
	ID        uint64    `json:"id" pg:",pk"`
	ChatID    uint64    `json:"chat_id" pg:",notnull"`
	UserID    uint64    `json:"user_id" pg:",notnull"`
	Content   string    `json:"content" pg:",notnull"`
	Timestamp time.Time `json:"timestamp" pg:",notnull"`
}
