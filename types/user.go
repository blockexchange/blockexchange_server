package types

type UserType string
type UserRole string

const (
	UserTypeGithub  UserType = "GITHUB"
	UserTypeLocal   UserType = "LOCAL"
	UserTypeDiscord UserType = "DISCORD"
	UserTypeMesehub UserType = "MESEHUB"
	UserTypeCDB     UserType = "CDB"
)

const (
	UserRoleDefault UserRole = "DEFAULT"
	UserRoleAdmin   UserRole = "ADMIN"
)

type User struct {
	ID         *int64   `json:"id"`
	Created    int64    `json:"created"`
	Name       string   `json:"name"`
	Hash       string   `json:"-"` // not exported
	Type       UserType `json:"type"`
	Role       UserRole `json:"role"`
	ExternalID *string  `json:"external_id"`
}

type UserSearch struct {
	Name   *string `json:"name"`
	Limit  *int    `json:"limit"`
	Offset *int    `json:"offset"`
}

func (u *User) Columns(action string) []string {
	cols := []string{}
	if action != "insert" {
		cols = append(cols, "id")
	}
	cols = append(cols, "created", "name", "hash", "type", "role", "external_id")
	return cols
}

func (u *User) Table() string {
	return "public.user"
}

func (u *User) Scan(action string, r func(dest ...any) error) error {
	return r(&u.ID, &u.Created, &u.Name, &u.Hash, &u.Type, &u.Role, &u.ExternalID)
}

func (u *User) Values(action string) []any {
	vals := []any{}
	if action != "insert" {
		vals = append(vals, u.ID)
	}
	vals = append(vals, u.Created, u.Name, u.Hash, u.Type, u.Role, u.ExternalID)
	return vals
}
