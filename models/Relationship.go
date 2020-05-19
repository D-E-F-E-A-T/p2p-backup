package models

type Relationship struct {
	ID      uint64 `json:"id" pg:",nopk,notnull"` // Target of relation
	UserID  uint64 `json:"user" pg:",notnull"`    // Owner of relation
	Blocked bool   `json:"status" pg:",use_zero"`
	Trusted bool   `json:"trusted" pg:",use_zero"`
}
