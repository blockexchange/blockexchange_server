package api_test

import (
	"blockexchange/api"
	"blockexchange/testutils"
	"blockexchange/types"
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func importRequest(t *testing.T, a *api.Api, user *types.User, filename string, body []byte) int {
	r := httptest.NewRequest("POST", "http://", bytes.NewReader(body))
	w := httptest.NewRecorder()
	r = mux.SetURLVars(r, map[string]string{"filename": filename})

	a.ImportSchematic(w, r, &api.SecureContext{
		Claims: &types.Claims{Username: user.Name, UserUID: user.UID},
	})
	return w.Result().StatusCode
}

func TestImportWorldedit(t *testing.T) {
	a := NewTestApi(t)
	user := testutils.CreateUser(a.UserRepo, t, nil)

	data, err := os.ReadFile("../schematic/worldedit/testdata/plain_chest.we")
	assert.NoError(t, err)

	assert.Equal(t, 200, importRequest(t, a, user, "chest.we", data))
}

func TestImportInvalidWorldedit(t *testing.T) {
	a := NewTestApi(t)
	user := testutils.CreateUser(a.UserRepo, t, nil)

	assert.Equal(t, 400, importRequest(t, a, user, "x.we", []byte("5:while true do end")))
	assert.Equal(t, 400, importRequest(t, a, user, "x.we", []byte("5")))
	assert.Equal(t, 400, importRequest(t, a, user, "x.txt", []byte("hello")))
}

func TestImportTooLarge(t *testing.T) {
	a := NewTestApi(t)
	user := testutils.CreateUser(a.UserRepo, t, nil)

	assert.Equal(t, 413, importRequest(t, a, user, "x.we", make([]byte, api.MAX_SIZE+1)))
}

func TestImportBX(t *testing.T) {
	a := NewTestApi(t)
	user := testutils.CreateUser(a.UserRepo, t, nil)

	data, err := os.ReadFile("../schematic/bx/testdata/blockexchange.zip")
	assert.NoError(t, err)

	assert.Equal(t, 200, importRequest(t, a, user, "schema.zip", data))
	assert.Equal(t, 400, importRequest(t, a, user, "schema.zip", []byte("not a zip")))
}

func TestBXExportImportDependencies(t *testing.T) {
	a := NewTestApi(t)
	owner := testutils.CreateUser(a.UserRepo, t, nil)
	other := testutils.CreateUser(a.UserRepo, t, nil)
	// a real schematic (the preview is rendered on import)
	importSchema := func(user *types.User, data []byte) *types.Schema {
		r := mux.SetURLVars(httptest.NewRequest("POST", "http://", bytes.NewReader(data)), map[string]string{"filename": "schema.zip"})
		w := httptest.NewRecorder()
		a.ImportSchematic(w, r, &api.SecureContext{Claims: &types.Claims{Username: user.Name, UserUID: user.UID}})
		assert.Equal(t, 200, w.Result().StatusCode, w.Body.String())
		schema := &types.Schema{}
		assert.NoError(t, json.Unmarshal(w.Body.Bytes(), schema))
		assert.NotEmpty(t, schema.UID)
		return schema
	}
	data, err := os.ReadFile("../schematic/bx/testdata/blockexchange.zip")
	assert.NoError(t, err)
	schema := importSchema(owner, data)
	deps, err := a.SchemaModDependencyRepo.GetBySchemaUID(schema.UID)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(deps))

	assert.NoError(t, a.SchemaModDependencyRepo.Create(&types.SchemaModDependency{
		SchemaUID: schema.UID, ModType: types.SchemaModDependencyTypeGame, GitURL: "https://github.com/luanti-org/minetest_game",
	}))
	assert.NoError(t, a.SchemaModDependencyRepo.Create(&types.SchemaModDependency{
		SchemaUID: schema.UID, ModType: types.SchemaModDependencyTypeMod, GitURL: "https://github.com/mt-mods/promise", Version: "v1.0",
	}))

	// export (public)
	r := mux.SetURLVars(httptest.NewRequest("GET", "http://", nil), map[string]string{"schema_uid": schema.UID})
	w := httptest.NewRecorder()
	a.ExportBXSchema(w, r)
	assert.Equal(t, 200, w.Result().StatusCode)

	// import by another user
	imported := importSchema(other, w.Body.Bytes())
	assert.NotEqual(t, schema.UID, imported.UID)

	deps, err = a.SchemaModDependencyRepo.GetBySchemaUID(imported.UID)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(deps))
	versions := map[string]string{}
	for _, d := range deps {
		assert.NotEmpty(t, d.UID)
		versions[d.GitURL] = string(d.ModType) + "@" + d.Version
	}
	assert.Equal(t, map[string]string{
		"https://github.com/luanti-org/minetest_game": "game@latest",
		"https://github.com/mt-mods/promise":          "mod@v1.0",
	}, versions)

	// the original dependencies are untouched
	deps, err = a.SchemaModDependencyRepo.GetBySchemaUID(schema.UID)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(deps))
}
