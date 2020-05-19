package models

import (
	"time"
)

type Statistic struct {
	ID              uint64    `json:"id" pg:",pk"`
	UserID          uint64    `json:"user_id" pg:",notnull"`
	KnownAddresses  []string  `json:"known_addresses"`
	EmailVerified   bool      `json:"email_verified" pg:",use_zero"`
	PhoneVerified   bool      `json:"phone_verified" pg:",use_zero"`
	DocsVerified    bool      `json:"docs_verified" pg:",use_zero"`
	LastSeen        time.Time `json:"last_seen"`
	Created         time.Time `json:"created"`
	FirstDeal       time.Time `json:"first_deal"`
	LastDeal        time.Time `json:"last_deal"`
	FinishedDeals   uint      `json:"finished_deals"`
	DealsVolume     uint      `json:"deals_volume"`
	Rating          uint      `json:"rating"`
	ActiveReferrals uint      `json:"active_referrals"`
}
