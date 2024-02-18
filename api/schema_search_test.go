package api_test

import (
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
