package types

import (
	"github.com/golang-jwt/jwt/v4"
)

type JWTPermission string

const (
	JWTPermissionUpload     JWTPermission = "UPLOAD"
	JWTPermissionOverwrite  JWTPermission = "OVERWRITE"
	JWTPermissionManagement JWTPermission = "MANAGEMENT"
	JWTPermissionMedia      JWTPermission = "MEDIA"
	JWTPermissionAdmin      JWTPermission = "ADMIN"
)

type Claims struct {
	*jwt.RegisteredClaims
	Username    string          `json:"username"`
	UserUID     string          `json:"user_uid"`
	Type        UserType        `json:"type"`
	Permissions []JWTPermission `json:"permissions"`
}

func (c *Claims) HasPermission(perm JWTPermission) bool {
	for _, p := range c.Permissions {
		if p == perm {
			return true
		}
	}
	return false
}

func (c *Claims) IsAdmin() bool {
	return c.HasPermission(JWTPermissionAdmin)
}
