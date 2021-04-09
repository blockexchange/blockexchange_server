package web

import (
	"blockexchange/core"
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

func TestSchemaCreateNoUser(t *testing.T) {
	db_ := testutils.CreateTestDatabase(t)
	api := NewApi(db_, core.NewNoOpCache())
	schema := &types.Schema{
		UserID:      -1,
		Name:        "my-schema",
		Description: "something",
		SizeX:       16,
		SizeY:       16,
		SizeZ:       16,
		License:     "CC0",
	}

	data, err := json.Marshal(schema)
	assert.NoError(t, err)
	r := httptest.NewRequest("GET", "http://", bytes.NewBuffer(data))
	w := httptest.NewRecorder()

	Secure(api.CreateSchema)(w, r)
	assert.Equal(t, 401, w.Result().StatusCode)
}

func TestSchemaCreateInvalidUser(t *testing.T) {
	db_ := testutils.CreateTestDatabase(t)
	api := NewApi(db_, core.NewNoOpCache())
	user := testutils.CreateUser(api.UserRepo, t, &types.User{})
	schema := &types.Schema{
		UserID:      -1,
		Name:        "my-schema",
		Description: "something",
		SizeX:       16,
		SizeY:       16,
		SizeZ:       16,
		License:     "CC0",
	}

	data, err := json.Marshal(schema)
	assert.NoError(t, err)
	r := httptest.NewRequest("GET", "http://", bytes.NewBuffer(data))
	w := httptest.NewRecorder()
	testutils.Login(t, r, user)

	Secure(api.CreateSchema)(w, r)
	assert.Equal(t, 200, w.Result().StatusCode)
}

func TestSchemaCreate(t *testing.T) {
	db_ := testutils.CreateTestDatabase(t)
	api := NewApi(db_, core.NewNoOpCache())
	user := testutils.CreateUser(api.UserRepo, t, &types.User{})
	schema := &types.Schema{
		UserID:      user.ID,
		Name:        "my-schema",
		Description: "something",
		SizeX:       16,
		SizeY:       16,
		SizeZ:       16,
		License:     "CC0",
	}

	// create

	data, err := json.Marshal(schema)
	assert.NoError(t, err)
	r := httptest.NewRequest("GET", "http://", bytes.NewBuffer(data))
	w := httptest.NewRecorder()
	testutils.Login(t, r, user)

	Secure(api.CreateSchema)(w, r)
	assert.Equal(t, 200, w.Result().StatusCode)

	err = json.NewDecoder(w.Body).Decode(&schema)
	assert.NoError(t, err)

	// update
	schema.Name = "something else"

	data, err = json.Marshal(schema)
	assert.NoError(t, err)
	r = httptest.NewRequest("GET", "http://", bytes.NewBuffer(data))
	r = mux.SetURLVars(r, map[string]string{"id": strconv.Itoa(int(schema.ID))})
	w = httptest.NewRecorder()
	testutils.Login(t, r, user)

	Secure(api.UpdateSchema)(w, r)
	assert.Equal(t, 200, w.Result().StatusCode)

	// update infos
	r = httptest.NewRequest("POST", "http://", nil)
	r = mux.SetURLVars(r, map[string]string{"id": strconv.Itoa(int(schema.ID))})
	w = httptest.NewRecorder()
	testutils.Login(t, r, user)

	Secure(api.UpdateSchemaInfo)(w, r)
	fmt.Println(w.Body.String())
	assert.Equal(t, 200, w.Result().StatusCode)

}

func TestSchemaCreateAndDownload(t *testing.T) {
	db_ := testutils.CreateTestDatabase(t)
	api := NewApi(db_, core.NewNoOpCache())
	user := testutils.CreateUser(api.UserRepo, t, &types.User{})
	schema := &types.Schema{
		UserID:      user.ID,
		Name:        "my-schema",
		Description: "something",
		SizeX:       16,
		SizeY:       16,
		SizeZ:       16,
		License:     "CC0",
	}

	// create

	data, err := json.Marshal(schema)
	assert.NoError(t, err)
	r := httptest.NewRequest("GET", "http://", bytes.NewBuffer(data))
	w := httptest.NewRecorder()
	testutils.Login(t, r, user)

	Secure(api.CreateSchema)(w, r)
	assert.Equal(t, 200, w.Result().StatusCode)

	err = json.NewDecoder(w.Body).Decode(&schema)
	assert.NoError(t, err)

	// get by id
	r = httptest.NewRequest("GET", "http://", nil)
	r = mux.SetURLVars(r, map[string]string{"id": strconv.Itoa(int(schema.ID))})
	w = httptest.NewRecorder()

	api.GetSchema(w, r)
	assert.Equal(t, 200, w.Result().StatusCode)

	schema2 := types.Schema{}
	err = json.NewDecoder(w.Body).Decode(&schema2)
	assert.NoError(t, err)
	assert.Equal(t, schema.Name, schema2.Name)
	assert.Equal(t, schema.Description, schema2.Description)
	assert.Equal(t, schema.Downloads, schema2.Downloads)
	assert.Equal(t, schema.License, schema2.License)
	assert.Equal(t, schema.TotalSize, schema2.TotalSize)

	// download by id

	r = httptest.NewRequest("GET", "http://", nil)
	r = mux.SetURLVars(r, map[string]string{"id": strconv.Itoa(int(schema.ID))})
	q := r.URL.Query()
	q.Add("download", "true")
	r.URL.RawQuery = q.Encode()

	w = httptest.NewRecorder()

	api.GetSchema(w, r)
	assert.Equal(t, 200, w.Result().StatusCode)

	schema2 = types.Schema{}
	err = json.NewDecoder(w.Body).Decode(&schema2)
	assert.NoError(t, err)
	assert.Equal(t, schema.Downloads+1, schema2.Downloads)

	// update

	schema2.Description = "another description"
	data, err = json.Marshal(schema2)
	assert.NoError(t, err)

	r = httptest.NewRequest("GET", "http://", bytes.NewBuffer(data))
	w = httptest.NewRecorder()
	r = mux.SetURLVars(r, map[string]string{"id": strconv.Itoa(int(schema.ID))})
	testutils.Login(t, r, user)

	Secure(api.UpdateSchema)(w, r)

	schema3, err := api.SchemaRepo.GetSchemaById(schema2.ID)
	assert.NoError(t, err)
	assert.NotNil(t, schema3)
	assert.Equal(t, schema2.Description, schema3.Description)

	// delete

	r = httptest.NewRequest("DELETE", "http://", nil)
	r = mux.SetURLVars(r, map[string]string{"id": strconv.Itoa(int(schema.ID))})
	w = httptest.NewRecorder()
	testutils.Login(t, r, user)

	Secure(api.DeleteSchema)(w, r)
	assert.Equal(t, 200, w.Result().StatusCode)

}
