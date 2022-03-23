package web

import (
	"blockexchange/testutils"
	"blockexchange/types"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
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
	api := NewTestApi(t)
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

func TestCountUsers(t *testing.T) {
	api := NewTestApi(t)
	testutils.CreateUser(api.UserRepo, t, &types.User{})

	count, err := api.UserRepo.CountUsers()
	assert.NoError(t, err)
	assert.NotEqual(t, 0, count)

	r := httptest.NewRequest("GET", "http://?limit=10&offset=1", nil)
	w := httptest.NewRecorder()
	api.GetUsers(w, r)
	assert.Equal(t, 200, w.Result().StatusCode)

	response := &types.PagedUsersResponse{}
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
	assert.Equal(t, 10, response.Limit)
	assert.Equal(t, 1, response.Offset)
	assert.NotEqual(t, 0, len(response.List))
	assert.NotEqual(t, 0, response.Total)
}
