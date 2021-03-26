package types

type UserType string

const (
	UserTypeGithub  UserType = "GITHUB"
	UserTypeLocal   UserType = "LOCAL"
	UserTypeDiscord UserType = "DISCORD"
	UserTypeMesehub UserType = "MESEHUB"
)

type User struct {
	ID         int64    `json:"id" db:"id"`
	Created    int64    `json:"created" db:"created"`
	Name       string   `json:"name" db:"name"`
	Hash       string   `db:"hash"`
	Type       UserType `json:"type" db:"type"`
	ExternalID *string  `json:"external_id" db:"external_id"`
	Mail       *string  `json:"mail" db:"mail"`
}
