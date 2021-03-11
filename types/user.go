package types

const (
	UserTypeGithub  = "GITHUB"
	UserTypeLocal   = "LOCAL"
	UserTypeDiscord = "DISCORD"
	UserTypeMesehub = "MESEHUB"
)

type User struct {
	ID         int64  `json:"id" db:"id"`
	Created    int64  `json:"created" db:"created"`
	Name       string `json:"name" db:"name"`
	Hash       string `json:"hash" db:"hash"`
	Type       string `json:"type" db:"type"`
	ExternalID string `json:"external_id" db:"external_id"`
	Mail       string `json:"mail" db:"mail"`
}
