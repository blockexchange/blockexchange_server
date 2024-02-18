package api_test

import (
	"blockexchange/testutils"
	"blockexchange/types"
	"encoding/json"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
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

	r := httptest.NewRequest("GET", "http://", nil)
	w := httptest.NewRecorder()

	r = mux.SetURLVars(r, map[string]string{
		"schema_name": schema.Name,
		"user_name":   user.Name,
	})

	api.SearchSchemaByNameAndUser(w, r)
	assert.Equal(t, 200, w.Result().StatusCode)

	search_result := &types.SchemaSearchResponse{}
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), search_result))
	assert.Equal(t, schema.UID, search_result.Schema.UID)
	assert.Equal(t, 0, search_result.Schema.Stars)
	assert.Equal(t, 1, len(search_result.Mods))
	assert.Equal(t, "default", search_result.Mods[0])

}
