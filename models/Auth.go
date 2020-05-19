package models

import (
	"time"
)

const (
	AuthRegistration AuthPurpose = iota
	AuthDeactivation
	AuthEnter
	AuthDeal
)

type AuthPurpose int

type Auth struct {
	ID      uint64      `json:"id" pg:",nopk,notnull"`
	Token   string      `json:"token" pg:",notnull"`
	Purpose AuthPurpose `json:"purpose"`
	Expires time.Time   `json:"expires" pg:",notnull"`
}
