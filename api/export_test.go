package api_test

import (
	"blockexchange/testutils"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestWorldEditExport(t *testing.T) {
	api := NewTestApi(t)
	user := testutils.CreateUser(api.UserRepo, t, nil)
	assert.NotEmpty(t, user.UID)
	schema := testutils.CreateSchema(api.SchemaRepo, t, user, nil)
	testutils.CreateSchemaPart(api.SchemaPartRepo, t, schema, nil)

	w, r := testutils.CreateGETRequest(t, map[string]string{
		"schema_uid": schema.UID,
	})

	api.ExportWorldeditSchema(w, r)

	assert.Equal(t, 200, w.Result().StatusCode)
	assert.Equal(t, "text/plain; charset=utf-8", w.Result().Header.Get("Content-Type"))
}

func TestBXExport(t *testing.T) {
	api := NewTestApi(t)
	user := testutils.CreateUser(api.UserRepo, t, nil)
	schema := testutils.CreateSchema(api.SchemaRepo, t, user, nil)
	testutils.CreateSchemaPart(api.SchemaPartRepo, t, schema, nil)

	r := httptest.NewRequest("GET", "http://", nil)
	w := httptest.NewRecorder()
	r = mux.SetURLVars(r, map[string]string{
		"schema_uid": schema.UID,
	})

	api.ExportBXSchema(w, r)

	assert.Equal(t, 200, w.Result().StatusCode)
	assert.Equal(t, "application/zip", w.Result().Header.Get("Content-Type"))
}
