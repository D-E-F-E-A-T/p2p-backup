package models

type User struct {
	ID          uint64 `json:"id" pg:",pk"`
	Login       string `json:"login" pg:",unique,notnull"`
	Email       string `json:"email" pg:",unique,notnull"`
	Phone       string `json:"phone" pg:",unique"`
	Telegram    int64  `json:"telegram" pg:",unique"`
	Password    string `json:"password" pg:",notnull"`
	Referral    string `json:"referral" pg:",unique"`
	Activation  string `json:"activation" pg:",unique"`
	Description string `json:"description"`
	Groups      []int  `json:"groups" pg:",array"`
}
