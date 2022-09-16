package api

import (
	"blockexchange/core"
	"blockexchange/testutils"
	"blockexchange/types"
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAccessTokenLogin(t *testing.T) {
	api := NewTestApi(t)

	user := &types.User{
		Type: types.UserTypeLocal,
	}
	testutils.CreateUser(api.UserRepo, t, user)

	token := &types.AccessToken{
		UserID:  *user.ID,
		Expires: (time.Now().Unix() + 300) * 1000,
		Created: time.Now().Unix() * 1000,
		Token:   "abcdef",
	}
	testutils.CreateAccessToken(api.AccessTokenRepo, t, token)

	login := types.Login{}
	login.Username = user.Name
	login.Token = token.Token
	data, err := json.Marshal(login)
	assert.NoError(t, err)
	assert.NotNil(t, data)

	r := httptest.NewRequest("GET", "http://", bytes.NewReader(data))
	w := httptest.NewRecorder()

	api.RequestToken(w, r)

	assert.NotNil(t, w.Body)
	info, err := core.ParseJWT(w.Body.String())
	assert.NoError(t, err)
	assert.NotNil(t, info)
	assert.Equal(t, user.Name, info.Username)
}

func TestInvalidAccessTokenLogin(t *testing.T) {
	api := NewTestApi(t)

	user := &types.User{
		Type: types.UserTypeLocal,
	}
	testutils.CreateUser(api.UserRepo, t, user)

	token := &types.AccessToken{
		UserID:  *user.ID,
		Expires: (time.Now().Unix() + 300) * 1000,
		Created: time.Now().Unix() * 1000,
		Token:   "abcdef",
	}
	testutils.CreateAccessToken(api.AccessTokenRepo, t, token)

	login := types.Login{}
	login.Username = user.Name
	login.Token = "invalid token"
	data, err := json.Marshal(login)
	assert.NoError(t, err)
	assert.NotNil(t, data)

	r := httptest.NewRequest("GET", "http://", bytes.NewReader(data))
	w := httptest.NewRecorder()

	api.RequestToken(w, r)

	assert.NotNil(t, w.Body)
	info, err := core.ParseJWT(w.Body.String())
	assert.Error(t, err)
	assert.Nil(t, info)
}
