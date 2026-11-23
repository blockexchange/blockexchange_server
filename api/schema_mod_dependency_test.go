package api_test

import (
	api_ "blockexchange/api"
	"blockexchange/core"
	"blockexchange/testutils"
	"blockexchange/types"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func loginWith(t *testing.T, r *http.Request, user *types.User, perms ...types.JWTPermission) {
	c := core.New(types.CreateConfig(), nil)
	token, err := c.CreateJWT(user, perms, time.Hour)
	assert.NoError(t, err)
	r.Header.Set("Authorization", token)
}

func TestSchemaModDependency(t *testing.T) {
	api := NewTestApi(t)
	owner := testutils.CreateUser(api.UserRepo, t, nil)
	other := testutils.CreateUser(api.UserRepo, t, nil)
	schema := testutils.CreateSchema(api.SchemaRepo, t, owner, nil)
	other_schema := testutils.CreateSchema(api.SchemaRepo, t, other, nil)

	call := func(handler api_.SecureHandler, method string, vars map[string]string, body any, user *types.User, perms ...types.JWTPermission) *httptest.ResponseRecorder {
		var buf *bytes.Buffer
		if body != nil {
			data, _ := json.Marshal(body)
			buf = bytes.NewBuffer(data)
		} else {
			buf = bytes.NewBuffer(nil)
		}
		r := mux.SetURLVars(httptest.NewRequest(method, "http://", buf), vars)
		w := httptest.NewRecorder()
		loginWith(t, r, user, perms...)
		api.Secure(handler)(w, r)
		return w
	}
	schema_vars := map[string]string{"schema_uid": schema.UID}
	list := func(user *types.User, perms ...types.JWTPermission) ([]*types.SchemaModDependency, int) {
		w := call(api.GetSchemaModDependencies, "GET", schema_vars, nil, user, perms...)
		result := []*types.SchemaModDependency{}
		if w.Result().StatusCode == 200 {
			assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &result))
		}
		return result, w.Result().StatusCode
	}

	// empty list
	deps, status := list(owner)
	assert.Equal(t, 200, status)
	assert.Equal(t, 0, len(deps))

	// owner creates, version defaults to "latest"
	w := call(api.CreateSchemaModDependency, "POST", schema_vars, map[string]any{"mod_type": "game", "git_url": " https://example.com/game.git "}, owner)
	assert.Equal(t, 200, w.Result().StatusCode)
	dep := &types.SchemaModDependency{}
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), dep))
	assert.NotEmpty(t, dep.UID)
	assert.Equal(t, schema.UID, dep.SchemaUID)
	assert.Equal(t, types.SchemaModDependencyTypeGame, dep.ModType)
	assert.Equal(t, "https://example.com/game.git", dep.GitURL)
	assert.Equal(t, "latest", dep.Version)

	// second dependency
	w = call(api.CreateSchemaModDependency, "POST", schema_vars, map[string]any{"mod_type": "mod", "git_url": "https://example.com/mod.git"}, owner)
	assert.Equal(t, 200, w.Result().StatusCode)

	deps, status = list(owner)
	assert.Equal(t, 200, status)
	assert.Equal(t, 2, len(deps))

	// validation
	assert.Equal(t, 400, call(api.CreateSchemaModDependency, "POST", schema_vars, map[string]any{"mod_type": "other", "git_url": "x"}, owner).Result().StatusCode)
	assert.Equal(t, 400, call(api.CreateSchemaModDependency, "POST", schema_vars, map[string]any{"mod_type": "mod", "git_url": ""}, owner).Result().StatusCode)
	assert.Equal(t, 400, call(api.CreateSchemaModDependency, "POST", schema_vars, map[string]any{"mod_type": "mod", "git_url": "file:///etc"}, owner).Result().StatusCode)
	assert.Equal(t, 400, call(api.CreateSchemaModDependency, "POST", schema_vars, map[string]any{"mod_type": "mod", "git_url": "https://example.com/x.git", "version": "--upload-pack=x"}, owner).Result().StatusCode)

	// other non-admin user can't read or write, admin can read
	_, status = list(other)
	assert.Equal(t, 403, status)
	_, status = list(other, types.JWTPermissionAdmin)
	assert.Equal(t, 200, status)
	assert.Equal(t, 403, call(api.CreateSchemaModDependency, "POST", schema_vars, map[string]any{"mod_type": "mod", "git_url": "x"}, other).Result().StatusCode)

	dep_vars := map[string]string{"schema_uid": schema.UID, "uid": dep.UID}
	update := map[string]any{"mod_type": "mod", "git_url": "https://example.com/mod2.git", "version": "v1"}
	assert.Equal(t, 403, call(api.UpdateSchemaModDependency, "PUT", dep_vars, update, other).Result().StatusCode)

	// a dependency can't be reached through another (owned) schema
	foreign_vars := map[string]string{"schema_uid": other_schema.UID, "uid": dep.UID}
	assert.Equal(t, 404, call(api.UpdateSchemaModDependency, "PUT", foreign_vars, update, other).Result().StatusCode)
	assert.Equal(t, 404, call(api.DeleteSchemaModDependency, "DELETE", foreign_vars, nil, other).Result().StatusCode)

	// admin may update
	w = call(api.UpdateSchemaModDependency, "PUT", dep_vars, update, other, types.JWTPermissionAdmin)
	assert.Equal(t, 200, w.Result().StatusCode)
	stored, err := api.SchemaModDependencyRepo.GetByUID(dep.UID)
	assert.NoError(t, err)
	assert.Equal(t, schema.UID, stored.SchemaUID)
	assert.Equal(t, types.SchemaModDependencyTypeMod, stored.ModType)
	assert.Equal(t, "v1", stored.Version)

	// other non-admin user can't delete, owner can
	assert.Equal(t, 403, call(api.DeleteSchemaModDependency, "DELETE", dep_vars, nil, other).Result().StatusCode)
	assert.Equal(t, 204, call(api.DeleteSchemaModDependency, "DELETE", dep_vars, nil, owner).Result().StatusCode)
	assert.Equal(t, 404, call(api.DeleteSchemaModDependency, "DELETE", dep_vars, nil, owner).Result().StatusCode)

	deps, _ = list(owner)
	assert.Equal(t, 1, len(deps))
}

func TestWasmPackAccess(t *testing.T) {
	api := NewTestApi(t)
	owner := testutils.CreateUser(api.UserRepo, t, nil)
	other := testutils.CreateUser(api.UserRepo, t, nil)
	schema := testutils.CreateSchema(api.SchemaRepo, t, owner, nil)

	get := func(user *types.User, perms ...types.JWTPermission) int {
		r := mux.SetURLVars(httptest.NewRequest("GET", "http://", nil), map[string]string{"schema_uid": schema.UID})
		w := httptest.NewRecorder()
		loginWith(t, r, user, perms...)
		api.Secure(api.GetWasmPack)(w, r)
		return w.Result().StatusCode
	}

	// not the owner
	assert.Equal(t, 403, get(other))
	// owner and admin pass the check, but there are no dependencies yet
	assert.Equal(t, 404, get(owner))
	assert.Equal(t, 404, get(other, types.JWTPermissionAdmin))

	// not logged in
	r := mux.SetURLVars(httptest.NewRequest("GET", "http://", nil), map[string]string{"schema_uid": schema.UID})
	w := httptest.NewRecorder()
	api.Secure(api.GetWasmPack)(w, r)
	assert.Equal(t, 401, w.Result().StatusCode)
}
