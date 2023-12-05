package api_test

import (
	"blockexchange/testutils"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestWorldEditExport(t *testing.T) {
	api := NewTestApi(t)
	user := testutils.CreateUser(api.UserRepo, t, nil)
	assert.NotNil(t, user.ID)
	schema := testutils.CreateSchema(api.SchemaRepo, t, user, nil)
	testutils.CreateSchemaPart(api.SchemaPartRepo, t, schema, nil)

	w, r := testutils.CreateGETRequest(t, map[string]string{
		"id": strconv.Itoa(int(schema.ID)),
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
		"id": strconv.Itoa(int(schema.ID)),
	})

	api.ExportBXSchema(w, r)

	assert.Equal(t, 200, w.Result().StatusCode)
	assert.Equal(t, "application/zip", w.Result().Header.Get("Content-Type"))
}
