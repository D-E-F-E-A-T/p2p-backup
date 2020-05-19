package models

import (
	"time"
)

type Notification struct {
	ID        uint64    `json:"id" pg:",nopk,notnull"`
	Title     string    `json:"title" pg:",notnull"`
	Link      string    `json:"link"`
	Content   string    `json:"content"`
	Seen      bool      `json:"seen"`
	Timestamp time.Time `json:"timestamp" pg:",notnull"`
}

type NotificationAgent interface {
}
