package models

import (
	"time"
)

type Comment struct {
	ID        uint64    `json:"id" pg:",nopk,notnull"`                // Target of comment
	UserID    uint64    `json:"user_id" pg:",notnull,unique:comment"` // Owner of comment
	Deal      uint64    `json:"deal" pg:",notnull,unique:comment"`
	Message   string    `json:"message" pg:",notnull"`
	Rating    uint8     `json:"rating" pg:",notnull"`
	Timestamp time.Time `json:"timestamp" pg:",notnull"`
}
