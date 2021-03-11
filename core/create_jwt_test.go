package core

import (
	"blockexchange/types"
	"fmt"
	"testing"

	jwt "github.com/dgrijalva/jwt-go"
)

func TestCreateJWT(t *testing.T) {
	user := types.User{
		Name: "dummy",
		ID:   123,
		Type: types.UserTypeGithub,
	}
	permissions := []string{"ADMIN"}
	token, err := CreateJWT(&user, permissions)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(token)

	parsedtoken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte{}, nil
	})

	if err != nil {
		t.Fatal(err)
	}
	if !parsedtoken.Valid {
		t.Fatal("invalid token")
	}
}
