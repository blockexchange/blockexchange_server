package core

import (
	"blockexchange/types"
	"errors"
	"fmt"
	"os"

	"github.com/golang-jwt/jwt/v4"
)

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

func CreateJWT(user *types.User, permissions []types.JWTPermission, exp int64) (string, error) {
	secret := os.Getenv("BLOCKEXCHANGE_KEY")
	claims := jwt.MapClaims{}
	claims["username"] = user.Name
	claims["user_id"] = user.ID
	claims["type"] = user.Type
	claims["mail"] = user.Mail
	claims["permissions"] = permissions
	claims["exp"] = exp
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), claims)
	return token.SignedString([]byte(secret))
}

func ParseJWT(token string) (*types.TokenInfo, error) {
	parsedtoken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("BLOCKEXCHANGE_KEY")), nil
	})

	if err != nil || !parsedtoken.Valid {
		return nil, err
	}

	info := types.TokenInfo{}
	claims, ok := parsedtoken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid claims object")
	}

	info.Username = claims["username"].(string)
	info.UserID = int64(claims["user_id"].(float64))
	info.Type = claims["type"].(string)
	if claims["mail"] != nil {
		info.Mail = claims["mail"].(string)
	}
	info.Permissions = make([]types.JWTPermission, 0)
	permissions := claims["permissions"].([]interface{})
	for _, permission := range permissions {
		info.Permissions = append(info.Permissions, types.JWTPermission(fmt.Sprint(permission)))
	}

	return &info, nil
}
