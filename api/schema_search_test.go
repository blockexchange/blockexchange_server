package api_test

import (
	"blockexchange/testutils"
	"blockexchange/types"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestSearchSchema(t *testing.T) {
	api := NewTestApi(t)

	// prepare data

	user := testutils.CreateUser(api.UserRepo, t, &types.User{})
	schema := types.Schema{
		UserID: *user.ID,
	}
	err := api.SchemaRepo.CreateSchema(&schema)
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

	search_result := &types.Schema{}
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), search_result))
	assert.Equal(t, *schema.ID, *search_result.ID)
	assert.Equal(t, 0, search_result.Stars)

}
