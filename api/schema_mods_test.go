package api_test

import (
	"blockexchange/testutils"
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestSchemaMods(t *testing.T) {
	api := NewTestApi(t)
	user := testutils.CreateUser(api.UserRepo, t, nil)
	schema := testutils.CreateSchema(api.SchemaRepo, t, user, nil)

	mods := []string{"x", "y"}
	buf, err := json.Marshal(&mods)
	assert.NoError(t, err)

	// Create
	r := httptest.NewRequest("POST", "http://", bytes.NewBuffer(buf))
	w := httptest.NewRecorder()
	Login(t, r, user)

	r = mux.SetURLVars(r, map[string]string{
		"schema_uid": schema.UID,
	})

	api.Secure(api.CreateSchemaMods)(w, r)

	assert.Equal(t, 204, w.Result().StatusCode)

	schema_mods, err := api.SchemaModRepo.GetSchemaModsBySchemaUID(schema.UID)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(schema_mods))

	// Get
	r = httptest.NewRequest("GET", "http://", nil)
	w = httptest.NewRecorder()
	Login(t, r, user)

	r = mux.SetURLVars(r, map[string]string{
		"schema_uid": schema.UID,
	})

	api.GetSchemaMods(w, r)

	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &mods))
	assert.NotNil(t, mods)
	assert.Equal(t, 2, len(mods))

}
