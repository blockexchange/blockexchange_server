package web

import (
	"blockexchange/core"
	"blockexchange/testutils"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestWorldEditExport(t *testing.T) {
	db_ := testutils.CreateTestDatabase(t)
	api := NewApi(db_, core.NewNoOpCache())
	user := testutils.CreateUser(api.UserRepo, t, nil)
	schema := testutils.CreateSchema(api.SchemaRepo, t, user, nil)
	testutils.CreateSchemaPart(api.SchemaPartRepo, t, schema, nil)

	r := httptest.NewRequest("GET", "http://", nil)
	w := httptest.NewRecorder()
	r = mux.SetURLVars(r, map[string]string{
		"username":   user.Name,
		"schemaname": schema.Name,
	})

	api.ExportWorldeditSchema(w, r)

	assert.Equal(t, 200, w.Result().StatusCode)
	assert.Equal(t, "text/plain; charset=utf-8", w.Result().Header.Get("Content-Type"))
}

func TestBXExport(t *testing.T) {
	db_ := testutils.CreateTestDatabase(t)
	api := NewApi(db_, core.NewNoOpCache())
	user := testutils.CreateUser(api.UserRepo, t, nil)
	schema := testutils.CreateSchema(api.SchemaRepo, t, user, nil)
	testutils.CreateSchemaPart(api.SchemaPartRepo, t, schema, nil)

	r := httptest.NewRequest("GET", "http://", nil)
	w := httptest.NewRecorder()
	r = mux.SetURLVars(r, map[string]string{
		"username":   user.Name,
		"schemaname": schema.Name,
	})

	api.ExportBXSchema(w, r)

	assert.Equal(t, 200, w.Result().StatusCode)
	assert.Equal(t, "application/zip", w.Result().Header.Get("Content-Type"))
}
