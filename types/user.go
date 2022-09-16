package types

type UserType string
type UserRole string

const (
	UserTypeGithub  UserType = "GITHUB"
	UserTypeLocal   UserType = "LOCAL"
	UserTypeDiscord UserType = "DISCORD"
	UserTypeMesehub UserType = "MESEHUB"
)

const (
	UserRoleDefault UserRole = "DEFAULT"
	UserRoleAdmin   UserRole = "ADMIN"
)

type User struct {
	ID         *int64 `json:"id"`
	Created    int64  `json:"created"`
	Name       string `json:"name"`
	Hash       string
	Type       UserType `json:"type"`
	Role       UserRole `json:"role"`
	ExternalID *string  `json:"external_id"`
	Mail       *string  `json:"mail"`
}

func (u *User) Columns(action string) []string {
	if action == "insert" {
		return []string{"created", "name", "hash", "type", "role", "external_id", "mail"}
	}
	return []string{"id", "created", "name", "hash", "type", "role", "external_id", "mail"}
}

func (u *User) Table() string {
	return "public.user"
}

func (u *User) Scan(action string, r func(dest ...any) error) error {
	return r(&u.ID, &u.Created, &u.Name, &u.Hash, &u.Type, &u.Role, &u.ExternalID, &u.Mail)
}

func (u *User) Values(action string) []any {
	if action == "insert" {
		return []any{u.Created, u.Name, u.Hash, u.Type, u.Role, u.ExternalID, u.Mail}
	}
	return []any{u.ID, u.Created, u.Name, u.Hash, u.Type, u.Role, u.ExternalID, u.Mail}
}
