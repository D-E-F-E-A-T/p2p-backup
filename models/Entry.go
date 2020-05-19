package models

import (
	"time"
)

type Entry struct {
	ID          uint64    `json:"id" pg:",nopk,notnull"`
	Address     string    `json:"address" pg:",notnull"`
	Application string    `json:"application" pg:",notnull"`
	Timestamp   time.Time `json:"timestamp" pg:",notnull"`
	Status      bool      `json:"status" pg:",use_zero"`
}
