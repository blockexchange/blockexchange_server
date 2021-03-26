package web

import (
	"blockexchange/core"
	"blockexchange/testutils"
	"blockexchange/types"
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

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

	data, err := json.Marshal(schema)
	assert.NoError(t, err)
	r := httptest.NewRequest("GET", "http://", bytes.NewBuffer(data))
	w := httptest.NewRecorder()
	testutils.Login(t, r, user)

	Secure(api.CreateSchema)(w, r)
	assert.Equal(t, 200, w.Result().StatusCode)
}
