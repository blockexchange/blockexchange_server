package web

import (
	"blockexchange/testutils"
	"blockexchange/types"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestAccessTokenNoAccess(t *testing.T) {
	api := NewTestApi(t)
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
	api := NewTestApi(t)
	user := testutils.CreateUser(api.UserRepo, t, &types.User{})
	token := &types.AccessToken{
		UserID:  user.ID,
		Expires: (time.Now().Unix() + 300) * 1000,
		Created: time.Now().Unix() * 1000,
		Token:   "abcdef",
	}

	data, err := json.Marshal(token)
	assert.NoError(t, err)
	r := httptest.NewRequest("GET", "http://", bytes.NewBuffer(data))
	w := httptest.NewRecorder()
	testutils.Login(t, r, user)

	Secure(api.PostAccessToken)(w, r)
	assert.Equal(t, 200, w.Result().StatusCode)

	r = httptest.NewRequest("GET", "http://", nil)
	w = httptest.NewRecorder()
	testutils.Login(t, r, user)

	// List tokens
	Secure(api.GetAccessTokens)(w, r)

	list := []types.AccessToken{}
	err = json.NewDecoder(w.Body).Decode(&list)
	assert.NoError(t, err)
	assert.NotNil(t, list)
	assert.Equal(t, 1, len(list))

	r, err = http.NewRequest("GET", "/foo", nil)
	assert.NoError(t, err)

	w = httptest.NewRecorder()
	testutils.Login(t, r, user)
	r = mux.SetURLVars(r, map[string]string{"id": strconv.Itoa(int(token.ID))})

	//Delete token
	Secure(api.DeleteAccessToken)(w, r)

	assert.Equal(t, 200, w.Result().StatusCode)

}
