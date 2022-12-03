package api

import (
	"blockexchange/testutils"
	"blockexchange/types"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestCreateSchemaPart(t *testing.T) {
	api := NewTestApi(t)
	user := testutils.CreateUser(api.UserRepo, t, nil)
	schema := testutils.CreateSchema(api.SchemaRepo, t, user, nil)
	schemapart := testutils.CreateSchemaPart(api.SchemaPartRepo, t, schema, nil)

	data, err := json.Marshal(schemapart)
	assert.NoError(t, err)

	// save

	r := httptest.NewRequest("GET", "http://", bytes.NewBuffer(data))
	w := httptest.NewRecorder()
	testutils.Login(t, r, user)

	Secure(api.CreateSchemaPart)(w, r)

	assert.Equal(t, 200, w.Result().StatusCode)

	// load

	r = httptest.NewRequest("GET", "http://", nil)
	w = httptest.NewRecorder()

	r = mux.SetURLVars(r, map[string]string{
		"schema_id": strconv.Itoa(int(schemapart.SchemaID)),
		"x":         strconv.Itoa(int(schemapart.OffsetX)),
		"y":         strconv.Itoa(int(schemapart.OffsetY)),
		"z":         strconv.Itoa(int(schemapart.OffsetZ)),
	})

	api.GetSchemaPart(w, r)

	assert.Equal(t, 200, w.Result().StatusCode)
	response_schemapart := &types.SchemaPart{}
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), response_schemapart))
	assert.Equal(t, schemapart.Data, response_schemapart.Data)
	assert.Equal(t, schemapart.MetaData, response_schemapart.MetaData)
}

func TestGetNextSchemaPart(t *testing.T) {
	api := NewTestApi(t)
	user := testutils.CreateUser(api.UserRepo, t, nil)
	schema := testutils.CreateSchema(api.SchemaRepo, t, user, nil)
	testutils.CreateSchemaPart(api.SchemaPartRepo, t, schema, &types.SchemaPart{
		OffsetX:  0,
		OffsetY:  0,
		OffsetZ:  0,
		Mtime:    100,
		SchemaID: schema.ID,
	})

	testutils.CreateSchemaPart(api.SchemaPartRepo, t, schema, &types.SchemaPart{
		OffsetX:  16,
		OffsetY:  0,
		OffsetZ:  0,
		Mtime:    200,
		SchemaID: schema.ID,
	})

	// load first

	r := httptest.NewRequest("GET", "http://", nil)
	w := httptest.NewRecorder()

	r = mux.SetURLVars(r, map[string]string{
		"schema_id": strconv.Itoa(int(schema.ID)),
	})

	api.GetFirstSchemaPart(w, r)

	assert.Equal(t, 200, w.Result().StatusCode)
	response_schemapart := &types.SchemaPart{}
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), response_schemapart))
	assert.Equal(t, 0, response_schemapart.OffsetX)
	assert.Equal(t, 0, response_schemapart.OffsetY)
	assert.Equal(t, 0, response_schemapart.OffsetZ)

	// load next

	r = httptest.NewRequest("GET", "http://", nil)
	w = httptest.NewRecorder()

	r = mux.SetURLVars(r, map[string]string{
		"schema_id": strconv.Itoa(int(schema.ID)),
		"x":         "0",
		"y":         "0",
		"z":         "0",
	})

	api.GetNextSchemaPart(w, r)

	assert.Equal(t, 200, w.Result().StatusCode)
	response_schemapart = &types.SchemaPart{}
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), response_schemapart))
	assert.Equal(t, 16, response_schemapart.OffsetX)
	assert.Equal(t, 0, response_schemapart.OffsetY)
	assert.Equal(t, 0, response_schemapart.OffsetZ)

	// load by mtime

	r = httptest.NewRequest("GET", "http://", nil)
	w = httptest.NewRecorder()

	r = mux.SetURLVars(r, map[string]string{
		"schema_id": fmt.Sprintf("%d", schema.ID),
		"mtime":     fmt.Sprintf("%d", 100),
	})

	api.GetNextSchemaPartByMtime(w, r)

	assert.Equal(t, 200, w.Result().StatusCode)
	response_schemapart = &types.SchemaPart{}
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), response_schemapart))
	assert.Equal(t, int64(200), response_schemapart.Mtime)
	assert.Equal(t, 16, response_schemapart.OffsetX)
	assert.Equal(t, 0, response_schemapart.OffsetY)
	assert.Equal(t, 0, response_schemapart.OffsetZ)

	// count by mtime

	r = httptest.NewRequest("GET", "http://", nil)
	w = httptest.NewRecorder()

	r = mux.SetURLVars(r, map[string]string{
		"schema_id": fmt.Sprintf("%d", schema.ID),
		"mtime":     fmt.Sprintf("%d", 100),
	})

	api.CountNextSchemaPartByMtime(w, r)

	assert.Equal(t, 200, w.Result().StatusCode)
	assert.Equal(t, []byte{'1'}, w.Body.Bytes())
}

func TestCreateSchemaPartInvalidSchemaID(t *testing.T) {
	api := NewTestApi(t)
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
