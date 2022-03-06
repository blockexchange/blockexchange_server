package web

import (
	"blockexchange/core"
	"blockexchange/testutils"
	"blockexchange/types"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestSchemaStar(t *testing.T) {
	db_ := testutils.CreateTestDatabase(t)
	api := NewApi(db_, core.NewNoOpCache())
	user := testutils.CreateUser(api.UserRepo, t, nil)
	schema := testutils.CreateSchema(api.SchemaRepo, t, user, nil)

	count, err := api.SchemaStarRepo.CountBySchemaID(schema.ID)
	assert.NoError(t, err)
	assert.Equal(t, 0, count)

	// Create
	r := httptest.NewRequest("GET", "http://", nil)
	w := httptest.NewRecorder()
	testutils.Login(t, r, user)

	r = mux.SetURLVars(r, map[string]string{
		"schema_id": strconv.Itoa(int(schema.ID)),
	})

	Secure(api.CreateSchemaStar)(w, r)

	assert.Equal(t, 200, w.Result().StatusCode)

	count, err = api.SchemaStarRepo.CountBySchemaID(schema.ID)
	assert.NoError(t, err)
	assert.Equal(t, 1, count)

	// GET without user_id
	r = httptest.NewRequest("GET", "http://", nil)
	r = mux.SetURLVars(r, map[string]string{
		"schema_id": strconv.Itoa(int(schema.ID)),
	})
	w = httptest.NewRecorder()
	api.GetSchemaStars(w, r)
	star_response := &types.SchemaStarResponse{}
	err = json.Unmarshal(w.Body.Bytes(), star_response)
	assert.NoError(t, err)
	assert.Equal(t, 1, star_response.Count)
	assert.Equal(t, false, star_response.Starred)

	// GET with user_id
	r = httptest.NewRequest("GET", fmt.Sprintf("http://?user_id=%d", user.ID), nil)
	r = mux.SetURLVars(r, map[string]string{
		"schema_id": strconv.Itoa(int(schema.ID)),
		"user_id":   strconv.Itoa(int(user.ID)),
	})
	w = httptest.NewRecorder()
	api.GetSchemaStars(w, r)
	star_response = &types.SchemaStarResponse{}
	err = json.Unmarshal(w.Body.Bytes(), star_response)
	assert.NoError(t, err)
	assert.Equal(t, 1, star_response.Count)
	assert.Equal(t, true, star_response.Starred)

	// Delete
	r = httptest.NewRequest("GET", "http://", nil)
	w = httptest.NewRecorder()
	testutils.Login(t, r, user)

	r = mux.SetURLVars(r, map[string]string{
		"schema_id": strconv.Itoa(int(schema.ID)),
	})

	Secure(api.DeleteSchemaStar)(w, r)

	count, err = api.SchemaStarRepo.CountBySchemaID(schema.ID)
	assert.NoError(t, err)
	assert.Equal(t, 0, count)

}
