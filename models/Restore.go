package models

import (
	"time"
)

type Restore struct {
	ID      uint64     `json:"id" pg:",pk"`
	Token   string    `json:"token" pg:",pk"`
	Expires time.Time `json:"expires" pg:",notnull"`
}
