package api_test

import (
	"blockexchange/testutils"
	"blockexchange/types"
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestSchemaCreateNoUser(t *testing.T) {
	api := NewTestApi(t)
	schema := &types.Schema{
		UserUID:     uuid.NewString(),
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

	api.Secure(api.CreateSchema)(w, r)
	assert.Equal(t, 401, w.Result().StatusCode)
}

func TestSchemaCreateInvalidUser(t *testing.T) {
	api := NewTestApi(t)
	user := testutils.CreateUser(api.UserRepo, t, &types.User{})
	schema := &types.Schema{
		UserUID:     uuid.NewString(),
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
	Login(t, r, user)

	api.Secure(api.CreateSchema)(w, r)
	assert.Equal(t, 200, w.Result().StatusCode)
}

func TestSchemaCreate(t *testing.T) {
	api := NewTestApi(t)
	user := testutils.CreateUser(api.UserRepo, t, &types.User{})
	schema := &types.Schema{
		UserUID:     uuid.NewString(),
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
	Login(t, r, user)

	api.Secure(api.CreateSchema)(w, r)
	assert.Equal(t, 200, w.Result().StatusCode)

	err = json.NewDecoder(w.Body).Decode(&schema)
	assert.NoError(t, err)

	// get by names
	schema2, err := api.SchemaRepo.GetSchemaByUsernameAndName(user.Name, schema.Name)
	assert.NoError(t, err)
	assert.NotNil(t, schema2)
	assert.Equal(t, schema.Description, schema2.Description)
	assert.Equal(t, *schema.ID, *schema2.ID)
	assert.Equal(t, schema.Created, schema2.Created)

	// update
	schema.Name = "something"

	data, err = json.Marshal(schema)
	assert.NoError(t, err)
	r = httptest.NewRequest("GET", "http://", bytes.NewBuffer(data))
	r = mux.SetURLVars(r, map[string]string{"schema_id": strconv.Itoa(int(*schema.ID))})
	w = httptest.NewRecorder()
	Login(t, r, user)

	api.Secure(api.UpdateSchema)(w, r)
	assert.Equal(t, 200, w.Result().StatusCode)

	// update infos
	r = httptest.NewRequest("POST", "http://", nil)
	r = mux.SetURLVars(r, map[string]string{"schema_id": strconv.Itoa(int(*schema.ID))})
	w = httptest.NewRecorder()
	Login(t, r, user)

	api.Secure(api.UpdateSchemaInfo)(w, r)
	assert.Equal(t, 200, w.Result().StatusCode)

}

func TestSchemaCreateAndDownload(t *testing.T) {
	api := NewTestApi(t)
	user := testutils.CreateUser(api.UserRepo, t, &types.User{})
	schema := &types.Schema{
		UserUID:     uuid.NewString(),
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
	Login(t, r, user)

	api.Secure(api.CreateSchema)(w, r)
	assert.Equal(t, 200, w.Result().StatusCode)

	err = json.NewDecoder(w.Body).Decode(&schema)
	assert.NoError(t, err)

	// get by id
	r = httptest.NewRequest("GET", "http://", nil)
	r = mux.SetURLVars(r, map[string]string{"schema_id": strconv.Itoa(int(*schema.ID))})
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
	r = mux.SetURLVars(r, map[string]string{"schema_id": strconv.Itoa(int(*schema.ID))})
	q := r.URL.Query()
	q.Add("download", "true")
	r.URL.RawQuery = q.Encode()

	w = httptest.NewRecorder()

	api.GetSchema(w, r)
	assert.Equal(t, 200, w.Result().StatusCode)

	schema2 = types.Schema{}
	err = json.NewDecoder(w.Body).Decode(&schema2)
	assert.NoError(t, err)

	// update

	schema2.Description = "another description"
	data, err = json.Marshal(schema2)
	assert.NoError(t, err)

	r = httptest.NewRequest("GET", "http://", bytes.NewBuffer(data))
	w = httptest.NewRecorder()
	r = mux.SetURLVars(r, map[string]string{"schema_id": strconv.Itoa(int(*schema.ID))})
	Login(t, r, user)

	api.Secure(api.UpdateSchema)(w, r)
	assert.Equal(t, 200, w.Result().StatusCode)

	schema3, err := api.SchemaRepo.GetSchemaById(*schema2.ID)
	assert.NoError(t, err)
	assert.NotNil(t, schema3)
	assert.Equal(t, schema2.Description, schema3.Description)

	// download by username and schemaname
	r = httptest.NewRequest("GET", "http://", bytes.NewBuffer(data))
	w = httptest.NewRecorder()
	r = mux.SetURLVars(r, map[string]string{
		"schema_name": schema.Name,
		"user_name":   user.Name,
	})
	q = r.URL.Query()
	q.Add("download", "true")
	r.URL.RawQuery = q.Encode()

	api.SearchSchemaByNameAndUser(w, r)
	assert.Equal(t, 200, w.Result().StatusCode)

	schema4 := types.Schema{}
	err = json.NewDecoder(w.Body).Decode(&schema4)
	assert.NoError(t, err)

	// counter updated
	schema3, err = api.SchemaRepo.GetSchemaById(*schema2.ID)
	assert.NoError(t, err)
	assert.NotNil(t, schema3)
}
