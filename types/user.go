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
	UID        string   `json:"uid" gorm:"primarykey;column:uid"`
	Created    int64    `json:"created" gorm:"column:created"`
	Name       string   `json:"name" gorm:"column:name"`
	Hash       string   `json:"-" gorm:"column:hash"` // not exported
	Type       UserType `json:"type" gorm:"column:type"`
	Role       UserRole `json:"role" gorm:"column:role"`
	ExternalID *string  `json:"external_id" gorm:"column:external_id"`
	AvatarURL  string   `json:"avatar_url" gorm:"column:avatar_url"`
}

func (u *User) TableName() string {
	return "user"
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
