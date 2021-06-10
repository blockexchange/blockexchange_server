package web

import (
	"blockexchange/core"
	"blockexchange/testutils"
	"blockexchange/types"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestSearchSchema(t *testing.T) {
	db_ := testutils.CreateTestDatabase(t)
	api := NewApi(db_, core.NewNoOpCache())

	user := testutils.CreateUser(api.UserRepo, t, &types.User{})
	schema := types.Schema{
		UserID: user.ID,
	}
	err := api.SchemaRepo.CreateSchema(&schema)
	assert.NoError(t, err)

	list, err := api.SchemaSearchRepo.FindByUsername(user.Name)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(list))
	assert.Equal(t, schema.ID, list[0].ID)

	r := httptest.NewRequest("GET", "http://", nil)
	w := httptest.NewRecorder()

	r = mux.SetURLVars(r, map[string]string{
		"schema_name": schema.Name,
		"user_name":   user.Name,
	})

	api.SearchSchemaByNameAndUser(w, r)

	assert.Equal(t, 200, w.Result().StatusCode)

	response_schema := &types.Schema{}
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), response_schema))
	assert.Equal(t, schema.ID, response_schema.ID)
}