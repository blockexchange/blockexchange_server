package web

import (
	"blockexchange/testutils"
	"blockexchange/types"
	"bytes"
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
		UserID: user.ID,
	}
	err := api.SchemaRepo.CreateSchema(&schema)
	assert.NoError(t, err)

	// search via search api

	order := types.CREATED
	order_dir := types.ASC
	complete := false
	q := &types.SchemaSearchRequest{
		UserName:       &user.Name,
		Complete:       &complete,
		OrderColumn:    &order,
		OrderDirection: &order_dir,
	}

	data, err := json.Marshal(q)
	assert.NoError(t, err)
	r := httptest.NewRequest("GET", "http://", bytes.NewBuffer(data))
	w := httptest.NewRecorder()
	api.SearchSchema(w, r)

	response_schema := &types.SchemaSearchResponse{}
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), response_schema))

	assert.Equal(t, 1, response_schema.Total)
	assert.Equal(t, 1, len(response_schema.List))
	assert.Equal(t, schema.ID, response_schema.List[0].ID)

	// search directly

	r = httptest.NewRequest("GET", "http://", nil)
	w = httptest.NewRecorder()

	r = mux.SetURLVars(r, map[string]string{
		"schema_name": schema.Name,
		"user_name":   user.Name,
	})

	api.SearchSchemaByNameAndUser(w, r)
	assert.Equal(t, 200, w.Result().StatusCode)

	search_result := &types.SchemaSearchResult{}
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), search_result))
	assert.Equal(t, schema.ID, search_result.ID)
	assert.Equal(t, 0, search_result.Stars)

	count, err := api.Repositories.SchemaSearchRepo.Count(q)
	assert.NoError(t, err)
	assert.Equal(t, 1, count)
}
