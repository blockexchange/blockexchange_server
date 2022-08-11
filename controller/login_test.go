package controller_test

import (
	"blockexchange/controller"
	"blockexchange/core"
	"blockexchange/db"
	"blockexchange/testutils"
	"blockexchange/types"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestLogin(t *testing.T) {
	db_ := testutils.CreateTestDatabase(t)
	repos := db.NewRepositories(db_)
	cfg := &core.Config{}
	c := controller.NewController(db_, cfg)

	// get

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "http://", nil)

	c.Login(w, r)

	assert.Equal(t, 200, w.Result().StatusCode)

	// post invalid login

	data := url.Values{}
	data.Set("action", "login")
	data.Set("username", "foo")
	data.Set("password", "bar")

	w = httptest.NewRecorder()
	r = httptest.NewRequest(http.MethodPost, "http://", strings.NewReader(data.Encode()))
	r.Header.Add("content-type", "application/x-www-form-urlencoded")

	c.Login(w, r)

	assert.Equal(t, 401, w.Result().StatusCode)

	// post valid login

	hash, err := bcrypt.GenerateFromPassword([]byte("bar"), bcrypt.DefaultCost)
	assert.NoError(t, err)
	user := testutils.CreateUser(repos.UserRepo, t, &types.User{
		Hash: string(hash),
	})

	data = url.Values{}
	data.Set("action", "login")
	data.Set("username", user.Name)
	data.Set("password", "bar")

	w = httptest.NewRecorder()
	r = httptest.NewRequest(http.MethodPost, "http://", strings.NewReader(data.Encode()))
	r.Header.Add("content-type", "application/x-www-form-urlencoded")

	c.Login(w, r)

	assert.Equal(t, 303, w.Result().StatusCode)

}
