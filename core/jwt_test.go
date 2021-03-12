package core

import (
	"blockexchange/types"
	"testing"
	"time"
)

func TestCreateJWT(t *testing.T) {
	user := types.User{
		Name: "dummy",
		ID:   123,
		Type: types.UserTypeGithub,
	}
	permissions := []types.JWTPermission{types.JWTPermissionUpload}
	exp := time.Now().Unix() + (3600 * 24 * 180)

	token, err := CreateJWT(&user, permissions, exp)
	if err != nil {
		t.Fatal(err)
	}

	info, err := ParseJWT(token)
	if err != nil {
		t.Fatal(err)
	}

	if info.Username != "dummy" {
		t.Fatal("username mismatch")
	}

	//TODO: check permissions, etc
}
