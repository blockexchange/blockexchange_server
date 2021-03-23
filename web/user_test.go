// +build integration

package web

import (
	"blockexchange/core"
	"blockexchange/db"
	"testing"

	"github.com/stretchr/testify/assert"
)

func assertValidUsername(t *testing.T, api *Api, username string) {
	valid, msg, err := api.ValidateUsername(username)
	if !valid || err != nil || msg != "" {
		t.Fail()
	}
}

func assertInvalidUsername(t *testing.T, api *Api, username string) {
	valid, msg, err := api.ValidateUsername(username)
	if valid || err != nil || msg == "" {
		t.Fail()
	}
}

func TestValidateUsername(t *testing.T) {
	db_, err := db.Init()
	assert.NoError(t, err)
	api := NewApi(db_, core.NewNoOpCache())

	assertValidUsername(t, api, "nonexistentuser")
	assertValidUsername(t, api, "SomeOne")
	assertValidUsername(t, api, "someone_else123-99")
	assertInvalidUsername(t, api, "")
	assertInvalidUsername(t, api, "invalid username")
	assertInvalidUsername(t, api, "invalid_username??")
	assertInvalidUsername(t, api, "invalid_username/")
	assertInvalidUsername(t, api, "invalid_usernameä")
	assertInvalidUsername(t, api, "invalid_username$")
}
