package web

import (
	"blockexchange/core"
	"blockexchange/testutils"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestGetStaticView(t *testing.T) {
	db_ := testutils.CreateTestDatabase(t)
	api := NewApi(db_, core.NewNoOpCache())

	r := httptest.NewRequest("GET", "http://", nil)
	w := httptest.NewRecorder()

	// no params, not found
	api.GetStaticView(w, r)
	assert.Equal(t, 404, w.Result().StatusCode)

	user := testutils.CreateUser(api.UserRepo, t, nil)
	schema := testutils.CreateSchema(api.SchemaRepo, t, user, nil)
	testutils.CreateSchemaScreenshot(api.SchemaScreenshotRepo, t, schema, nil)

	w = httptest.NewRecorder()
	r = mux.SetURLVars(r, map[string]string{
		"schema_name": schema.Name,
		"user_name":   user.Name,
	})
	api.GetStaticView(w, r)
	assert.Equal(t, 200, w.Result().StatusCode)
}
