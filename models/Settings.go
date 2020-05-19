package models

type Settings struct {
	ID                     uint64   `json:"id" pg:",pk"`
	UserID                 uint64   `json:"user_id" pg:",notnull"`
	Language               string   `json:"language" pg:"type:varchar(2)"`
	AuthenticatorKey       string   `json:"authenticator_key" pg:",unique"`
	AvailableAddresses     []string `json:"available_addresses" pg:",array"`
	TelegramAuthentication bool     `json:"telegram_authentication" pg:",use_zero"`
	Timezone               string   `json:"timezone"`
}
