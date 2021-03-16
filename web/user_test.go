package web

import (
	"blockexchange/types"
	"blockexchange/web/testdata"
	"net/http/httptest"
	"testing"
)

func TestGetUsers(t *testing.T) {
	r := httptest.NewRequest("GET", "http://", nil)
	w := httptest.NewRecorder()

	repo := testdata.MockUserRepository{}
	api := Api{UserRepo: &repo}

	api.GetUsers(w, r)
	//TODO: stub
}

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
	repo := testdata.MockUserRepository{
		Users: []types.User{
			types.User{
				Name: "somebody",
			},
		},
	}
	api := Api{UserRepo: &repo}

	assertValidUsername(t, &api, "nonexistentuser")
	assertValidUsername(t, &api, "SomeOne")
	assertValidUsername(t, &api, "someone_else123-99")
	assertInvalidUsername(t, &api, "")
	assertInvalidUsername(t, &api, "somebody")
	assertInvalidUsername(t, &api, "invalid username")
	assertInvalidUsername(t, &api, "invalid_username??")
	assertInvalidUsername(t, &api, "invalid_username/")
	assertInvalidUsername(t, &api, "invalid_usernameä")
	assertInvalidUsername(t, &api, "invalid_username$")
}
