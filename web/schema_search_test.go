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

	order := types.CREATED
	order_dir := types.ASC
	complete := false

	q := &types.SchemaSearchRequest{
		UserName:       &user.Name,
		Complete:       &complete,
		OrderColumn:    &order,
		OrderDirection: &order_dir,
	}
	list, err := api.SchemaSearchRepo.Search(q, 10, 0)
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

	response_schema := &types.SchemaSearchResult{}
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), response_schema))
	assert.Equal(t, schema.ID, response_schema.ID)
	assert.Equal(t, 0, response_schema.Stars)

	count, err := api.Repositories.SchemaSearchRepo.Count(q)
	assert.NoError(t, err)
	assert.Equal(t, 1, count)
}
