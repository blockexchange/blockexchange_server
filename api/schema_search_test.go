package api_test

import (
	"blockexchange/api"
	"blockexchange/testutils"
	"blockexchange/types"
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSearchSchema(t *testing.T) {
	api := NewTestApi(t)

	// prepare data

	user := testutils.CreateUser(api.UserRepo, t, &types.User{})
	schema := types.Schema{
		UserUID: user.UID,
		Mtime:   time.Now().Unix(),
		Created: time.Now().Unix(),
	}
	err := api.SchemaRepo.CreateSchema(&schema)
	api.SchemaModRepo.CreateSchemaMod(&types.SchemaMod{
		SchemaUID: schema.UID,
		ModName:   "default",
	})

	assert.NoError(t, err)

	// search directly

	data, err := json.Marshal(&types.SchemaSearchRequest{
		UserName:   &user.Name,
		SchemaName: &schema.Name,
	})
	assert.NoError(t, err)

	r := httptest.NewRequest("POST", "http://", bytes.NewBuffer(data))
	w := httptest.NewRecorder()

	api.SearchSchema(w, r)
	assert.Equal(t, 200, w.Result().StatusCode)

	search_result := []*types.SchemaSearchResponse{}
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &search_result))
	assert.Equal(t, 1, len(search_result))
	assert.Equal(t, schema.UID, search_result[0].Schema.UID)
	assert.Equal(t, 0, search_result[0].Schema.Stars)
	assert.Equal(t, 1, len(search_result[0].Mods))
	assert.Equal(t, "default", search_result[0].Mods[0])

}

func TestSearchStarredSchema(t *testing.T) {
	a := NewTestApi(t)

	owner := testutils.CreateUser(a.UserRepo, t, &types.User{})
	starrer := testutils.CreateUser(a.UserRepo, t, &types.User{})
	other := testutils.CreateUser(a.UserRepo, t, &types.User{})

	starred := types.Schema{UserUID: owner.UID, Mtime: time.Now().Unix(), Created: time.Now().Unix(), Complete: true}
	assert.NoError(t, a.SchemaRepo.CreateSchema(&starred))
	unstarred := types.Schema{UserUID: owner.UID, Mtime: time.Now().Unix(), Created: time.Now().Unix(), Complete: true}
	assert.NoError(t, a.SchemaRepo.CreateSchema(&unstarred))

	assert.NoError(t, a.SchemaStarRepo.Create(&types.SchemaStar{SchemaUID: starred.UID, UserUID: starrer.UID}))

	search := func(user *types.User) []*types.SchemaSearchResponse {
		r := httptest.NewRequest("POST", "http://", bytes.NewBufferString(`{"user_uid":"`+owner.UID+`"}`))
		w := httptest.NewRecorder()
		a.SearchStarredSchema(w, r, &api.SecureContext{Claims: &types.Claims{Username: user.Name, UserUID: user.UID}})
		assert.Equal(t, 200, w.Result().StatusCode)
		list := []*types.SchemaSearchResponse{}
		assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &list))
		return list
	}

	count := func(user *types.User) int64 {
		r := httptest.NewRequest("POST", "http://", bytes.NewBufferString(`{}`))
		w := httptest.NewRecorder()
		a.CountStarredSchema(w, r, &api.SecureContext{Claims: &types.Claims{Username: user.Name, UserUID: user.UID}})
		assert.Equal(t, 200, w.Result().StatusCode)
		var c int64
		assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &c))
		return c
	}

	list := search(starrer)
	assert.Equal(t, 1, len(list))
	assert.Equal(t, starred.UID, list[0].Schema.UID)
	assert.Equal(t, int64(1), count(starrer))

	// other users only see their own (empty) stars
	assert.Equal(t, 0, len(search(other)))
	assert.Equal(t, int64(0), count(other))

	// the star filter can't be injected through the public search endpoint
	r := httptest.NewRequest("POST", "http://", bytes.NewBufferString(`{"user_uid":"`+owner.UID+`","starred_by_user_uid":"`+starrer.UID+`","StarredByUserUID":"`+starrer.UID+`"}`))
	w := httptest.NewRecorder()
	a.SearchSchema(w, r)
	assert.Equal(t, 200, w.Result().StatusCode)
	public_list := []*types.SchemaSearchResponse{}
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &public_list))
	assert.Equal(t, 2, len(public_list))
}
