package core

import (
	"blockexchange/types"
	"os"

	jwt "github.com/dgrijalva/jwt-go"
)

func CreateJWT(user *types.User, permissions []string) (string, error) {
	secret := os.Getenv("BLOCKEXCHANGE_KEY")
	claims := jwt.MapClaims{}
	claims["username"] = user.Name
	claims["user_id"] = user.ID
	claims["type"] = user.Type
	claims["mail"] = user.Mail
	claims["permissions"] = permissions
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), claims)
	return token.SignedString([]byte(secret))
}
