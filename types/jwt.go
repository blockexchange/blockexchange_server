package types

import "fmt"

type JWTPermission string

const (
	JWTPermissionUpload     JWTPermission = "UPLOAD"
	JWTPermissionOverwrite                = "OVERWRITE"
	JWTPermissionManagement               = "MANAGEMENT"
)

type TokenInfo struct {
	Username    string
	ID          int
	Mail        string
	Type        string
	Permissions []JWTPermission
}

func (t TokenInfo) String() string {
	return fmt.Sprintf("TokenInfo: Username=%s ID=%d Mail=%s Type=%s Permissions=%s",
		t.Username, t.ID, t.Mail, t.Type, t.Permissions)
}
