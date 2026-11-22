package api_test

import (
	"blockexchange/core"
	"blockexchange/testutils"
	"blockexchange/types"
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestCreateAndUpdateTag(t *testing.T) {
	api := NewTestApi(t)
	user := testutils.CreateUser(api.UserRepo, t, nil)

	// create
	buf, _ := json.Marshal(map[string]any{"name": " my-tag ", "description": "desc", "restricted": true})
	r := httptest.NewRequest("POST", "http://", bytes.NewBuffer(buf))
	w := httptest.NewRecorder()
	Login(t, r, user)
	api.Secure(api.CreateTag)(w, r)
	assert.Equal(t, 200, w.Result().StatusCode)

	tag := &types.Tag{}
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), tag))
	assert.NotEmpty(t, tag.UID)
	assert.Equal(t, "my-tag", tag.Name)
	assert.True(t, tag.Restricted)

	// duplicate name
	r = httptest.NewRequest("POST", "http://", bytes.NewBuffer(buf))
	w = httptest.NewRecorder()
	Login(t, r, user)
	api.Secure(api.CreateTag)(w, r)
	assert.Equal(t, 409, w.Result().StatusCode)

	// empty name
	buf, _ = json.Marshal(map[string]any{"name": "  "})
	r = httptest.NewRequest("POST", "http://", bytes.NewBuffer(buf))
	w = httptest.NewRecorder()
	Login(t, r, user)
	api.Secure(api.CreateTag)(w, r)
	assert.Equal(t, 400, w.Result().StatusCode)

	// update, clearing restricted and description
	buf, _ = json.Marshal(map[string]any{"name": "renamed", "description": "", "restricted": false})
	r = httptest.NewRequest("PUT", "http://", bytes.NewBuffer(buf))
	r = mux.SetURLVars(r, map[string]string{"tag_uid": tag.UID})
	w = httptest.NewRecorder()
	Login(t, r, user)
	api.Secure(api.UpdateTag)(w, r)
	assert.Equal(t, 200, w.Result().StatusCode)

	stored, err := api.TagRepo.GetByUID(tag.UID)
	assert.NoError(t, err)
	assert.Equal(t, "renamed", stored.Name)
	assert.Equal(t, "", stored.Description)
	assert.False(t, stored.Restricted)

	// unknown tag
	r = httptest.NewRequest("PUT", "http://", bytes.NewBuffer(buf))
	r = mux.SetURLVars(r, map[string]string{"tag_uid": "00000000-0000-0000-0000-000000000000"})
	w = httptest.NewRecorder()
	Login(t, r, user)
	api.Secure(api.UpdateTag)(w, r)
	assert.Equal(t, 404, w.Result().StatusCode)
}

func TestCreateTagRequiresAdmin(t *testing.T) {
	api := NewTestApi(t)
	user := testutils.CreateUser(api.UserRepo, t, nil)

	c := core.New(types.CreateConfig(), nil)
	token, err := c.CreateJWT(user, []types.JWTPermission{types.JWTPermissionUpload}, time.Hour)
	assert.NoError(t, err)

	buf, _ := json.Marshal(map[string]any{"name": "x"})
	r := httptest.NewRequest("POST", "http://", bytes.NewBuffer(buf))
	r.Header.Set("Authorization", token)
	w := httptest.NewRecorder()
	api.Secure(api.CreateTag)(w, r)
	assert.Equal(t, 403, w.Result().StatusCode)
}
