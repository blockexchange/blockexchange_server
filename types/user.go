package types

type UserType string
type UserRole string

const (
	UserTypeLocal UserType = "LOCAL"
)

const (
	UserRoleDefault UserRole = "DEFAULT"
	UserRoleAdmin   UserRole = "ADMIN"
)

type User struct {
	UID        string   `json:"uid" ksql:"uid"`
	Created    int64    `json:"created" ksql:"created"`
	Name       string   `json:"name" ksql:"name"`
	Hash       string   `json:"-" ksql:"hash"` // not exported
	Type       UserType `json:"type" ksql:"type"`
	Role       UserRole `json:"role" ksql:"role"`
	ExternalID *string  `json:"external_id" ksql:"external_id"`
	AvatarURL  string   `json:"avatar_url" ksql:"avatar_url"`
}

type UserSearch struct {
	Name   *string `json:"name"`
	Limit  *int    `json:"limit"`
	Offset *int    `json:"offset"`
}

type ChangePassword struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}
