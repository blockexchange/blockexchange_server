package web

import (
	"blockexchange/core"
	"blockexchange/testutils"
	"blockexchange/types"
	"testing"
)

func assertValidUsername(t *testing.T, api *Api, username string) {
	valid, msg, err := api.ValidateUsername(username)
	if !valid || err != nil || msg != "" {
		t.Fatal(err)
	}
}

func assertInvalidUsername(t *testing.T, api *Api, username string) {
	valid, msg, err := api.ValidateUsername(username)
	if valid || err != nil || msg == "" {
		t.Fatal(err)
	}
}

func TestValidateUsername(t *testing.T) {
	db_ := testutils.CreateTestDatabase(t)
	api := NewApi(db_, core.NewNoOpCache())
	user := testutils.CreateUser(api.UserRepo, t, &types.User{})

	assertValidUsername(t, api, "nonexistentuser")
	assertValidUsername(t, api, "SomeOne")
	assertValidUsername(t, api, "someone_else123-99")
	assertInvalidUsername(t, api, "")
	assertInvalidUsername(t, api, user.Name)
	assertInvalidUsername(t, api, "invalid username")
	assertInvalidUsername(t, api, "invalid_username??")
	assertInvalidUsername(t, api, "invalid_username/")
	assertInvalidUsername(t, api, "invalid_usernameä")
	assertInvalidUsername(t, api, "invalid_username$")
}
