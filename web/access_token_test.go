package web

import (
	"blockexchange/core"
	"blockexchange/testutils"
	"blockexchange/types"
	"encoding/json"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAccessTokenNoAccess(t *testing.T) {
	db_ := testutils.CreateTestDatabase(t)
	api := NewApi(db_, core.NewNoOpCache())
	user := testutils.CreateUser(api.UserRepo, t, &types.User{})
	token := &types.AccessToken{
		UserID:  user.ID,
		Expires: (time.Now().Unix() + 300) * 1000,
		Created: time.Now().Unix() * 1000,
		Token:   "abcdef",
	}
	testutils.CreateAccessToken(api.AccessTokenRepo, t, token)

	r := httptest.NewRequest("GET", "http://", nil)
	w := httptest.NewRecorder()

	Secure(api.GetAccessTokens)(w, r)

	assert.Equal(t, 401, w.Result().StatusCode)
}

func TestAccessTokenFullAccess(t *testing.T) {
	db_ := testutils.CreateTestDatabase(t)
	api := NewApi(db_, core.NewNoOpCache())
	user := testutils.CreateUser(api.UserRepo, t, &types.User{})
	token := &types.AccessToken{
		UserID:  user.ID,
		Expires: (time.Now().Unix() + 300) * 1000,
		Created: time.Now().Unix() * 1000,
		Token:   "abcdef",
	}
	testutils.CreateAccessToken(api.AccessTokenRepo, t, token)

	r := httptest.NewRequest("GET", "http://", nil)
	w := httptest.NewRecorder()

	testutils.Login(t, r, user)

	Secure(api.GetAccessTokens)(w, r)

	list := []types.AccessToken{}
	err := json.NewDecoder(w.Body).Decode(&list)
	assert.NoError(t, err)
	assert.NotNil(t, list)
	assert.Equal(t, 1, len(list))
}
