package core

import (
	"blockexchange/types"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var err_unauthorized = errors.New("unauthorized")

func GetPermissions(user *types.User, management bool) []types.JWTPermission {
	permissions := []types.JWTPermission{
		types.JWTPermissionUpload,
		types.JWTPermissionOverwrite,
	}

	if management {
		permissions = append(permissions, types.JWTPermissionManagement)
	}

	if user.Role == types.UserRoleAdmin && management {
		permissions = append(permissions, types.JWTPermissionAdmin)
	}

	return permissions
}

func CreateJWT(user *types.User, permissions []types.JWTPermission, d time.Duration) (string, error) {
	c := types.Claims{
		UserID:      *user.ID,
		Username:    user.Name,
		Type:        user.Type,
		Permissions: permissions,
	}

	if user.Mail != nil {
		c.Mail = *user.Mail
	}

	c.RegisteredClaims = &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(d)),
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, c)

	return t.SignedString([]byte(os.Getenv("BLOCKEXCHANGE_KEY")))

}

func ParseJWT(token string) (*types.Claims, error) {
	t, err := jwt.ParseWithClaims(token, &types.Claims{}, func(token *jwt.Token) (any, error) {
		return []byte(os.Getenv("BLOCKEXCHANGE_KEY")), nil
	})

	if err != nil {
		return nil, err
	}

	if !t.Valid {
		return nil, err_unauthorized
	}

	claims, ok := t.Claims.(*types.Claims)
	if !ok {
		return nil, errors.New("internal error")
	}

	return claims, nil
}
