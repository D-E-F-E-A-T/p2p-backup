package models

import (
	"time"
)

type Offer struct {
	ID           uint64         `json:"id" pg:",pk"`
	UserID       uint64         `json:"user" pg:",notnull"`
	StatisticID  uint64         `json:"statistic_id" pg:",notnull"`
	SettingsID   uint64         `json:"settings_id" pg:",notnull"`
	WantCurrency uint8          `json:"want_currency" pg:",notnull"`
	HaveCurrency uint8          `json:"have_currency" pg:",notnull"`
	Minimal      float32        `json:"minimal"`
	Maximal      float32        `json:"maximal"`
	Profit       int16          `json:"profit"`
	Cost         float32        `json:"cost"`
	Title        string         `json:"title" pg:"type:varchar(40)"`
	Location     string         `json:"location" pg:"type:varchar(2)"`
	Provider     uint8          `json:"provider"`
	Terms        string         `json:"terms" pg:"type:varchar(256)"`
	Warranty     bool           `json:"warranty" pg:",use_zero"`
	WithDynamic  bool           `json:"with_dynamic" pg:",use_zero"`
	WithPhone    bool           `json:"with_phone" pg:",use_zero"`
	WithDocs     bool           `json:"with_docs" pg:",use_zero"`
	WithDeals    uint           `json:"with_deals"`
	WithRating   uint           `json:"with_rating"`
	Active       bool           `json:"active" pg:",use_zero"`
	Rating       uint32         `json:"rating"`
	CancelTime   uint8          `json:"cancel_time"`
	Days         []time.Weekday `json:"days" pg:",array"`
	StartTime    time.Time      `json:"start_time"`
	FinishTime   time.Time      `json:"finish_time"`
	Timestamp    time.Time      `json:"timestamp" pg:",notnull"`
}
