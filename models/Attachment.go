package models

type Attachment struct {
	MessageID uint64 `json:"message_id"`
	Path      string `json:"path"`
	Name      string `json:"name"`
}
