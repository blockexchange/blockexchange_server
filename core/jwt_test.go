package core

import (
	"blockexchange/types"
	"testing"
)

func TestCreateJWT(t *testing.T) {
	user := types.User{
		Name: "dummy",
		ID:   123,
		Type: types.UserTypeGithub,
	}
	permissions := []types.JWTPermission{types.JWTPermissionUpload}
	token, err := CreateJWT(&user, permissions)
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
