package web

import (
	"blockexchange/core"
	"blockexchange/testutils"
	"blockexchange/types"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestSchemaTag(t *testing.T) {
	db_ := testutils.CreateTestDatabase(t)
	api := NewApi(db_, core.NewNoOpCache())

	user := testutils.CreateUser(api.UserRepo, t, &types.User{})
	schema := types.Schema{
		UserID: user.ID,
	}
	err := api.SchemaRepo.CreateSchema(&schema)
	assert.NoError(t, err)

	tag := types.Tag{
		Name:        "mytag",
		Description: "desc",
	}
	err = api.TagRepo.Create(&tag)
	assert.NoError(t, err)

	r := httptest.NewRequest("GET", "http://", nil)
	w := httptest.NewRecorder()
	testutils.Login(t, r, user)

	r = mux.SetURLVars(r, map[string]string{
		"schema_id": strconv.Itoa(int(schema.ID)),
		"tag_id":    strconv.Itoa(int(tag.ID)),
	})

	Secure(api.CreateSchemaTag)(w, r)
	assert.Equal(t, 200, w.Result().StatusCode)

	Secure(api.DeleteSchemaTag)(w, r)
	assert.Equal(t, 200, w.Result().StatusCode)

	api.GetSchemaTags(w, r)
	assert.Equal(t, 200, w.Result().StatusCode)

}
