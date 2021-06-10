package web

import (
	"blockexchange/core"
	"blockexchange/testutils"
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
