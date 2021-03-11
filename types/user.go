package types

const (
	UserTypeGithub  = "GITHUB"
	UserTypeLocal   = "LOCAL"
	UserTypeDiscord = "DISCORD"
	UserTypeMesehub = "MESEHUB"
)

type User struct {
	ID         int64  `json:"id"`
	Created    int64  `json:"created"`
	Name       string `json:"name"`
	Hash       string `json:"hash"`
	Type       string `json:"type"`
	ExternalID string `json:"external_id"`
	Mail       string `json:"mail"`
}
