package models

import (
	"time"
)

type Session struct {
	ID      uint64     `json:"id" pg:",notnull"`
	Token   string    `json:"token" pg:",pk"`
	Expires time.Time `json:"expires" pg:",notnull"`
}
