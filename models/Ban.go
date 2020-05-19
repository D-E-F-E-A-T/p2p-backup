package models

import (
	"time"
)

type Ban struct {
	UserID      uint64    `json:"user_id" pg:",pk"`
	Expires     time.Time `json:"expires"`
	Description string    `json:"description"`
}
