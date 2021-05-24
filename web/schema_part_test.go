package web

import (
	"blockexchange/core"
	"blockexchange/testutils"
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateSchemaPart(t *testing.T) {
	db_ := testutils.CreateTestDatabase(t)
	api := NewApi(db_, core.NewNoOpCache())
	user := testutils.CreateUser(api.UserRepo, t, nil)
	schema := testutils.CreateSchema(api.SchemaRepo, t, user, nil)
	schemapart := testutils.CreateSchemaPart(api.SchemaPartRepo, t, schema, nil)

	data, err := json.Marshal(schemapart)
	assert.NoError(t, err)

	r := httptest.NewRequest("GET", "http://", bytes.NewBuffer(data))
	w := httptest.NewRecorder()
	testutils.Login(t, r, user)

	Secure(api.CreateSchemaPart)(w, r)

	assert.Equal(t, 200, w.Result().StatusCode)
}

func TestCreateSchemaPartInvalidSchemaID(t *testing.T) {
	db_ := testutils.CreateTestDatabase(t)
	api := NewApi(db_, core.NewNoOpCache())
	user := testutils.CreateUser(api.UserRepo, t, nil)
	schema := testutils.CreateSchema(api.SchemaRepo, t, user, nil)
	schemapart := testutils.CreateSchemaPart(api.SchemaPartRepo, t, schema, nil)
	schemapart.SchemaID = -1

	data, err := json.Marshal(schemapart)
	assert.NoError(t, err)

	r := httptest.NewRequest("GET", "http://", bytes.NewBuffer(data))
	w := httptest.NewRecorder()
	testutils.Login(t, r, user)

	Secure(api.CreateSchemaPart)(w, r)

	assert.Equal(t, 500, w.Result().StatusCode)
}
