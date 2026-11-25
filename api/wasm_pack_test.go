package api_test

import (
	"blockexchange/api"
	"blockexchange/core"
	"blockexchange/testutils"
	"blockexchange/types"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

// local repository with a single commit
func createRepo(t *testing.T, dir string, files map[string]string) {
	assert.NoError(t, os.MkdirAll(dir, 0755))
	for name, content := range files {
		assert.NoError(t, os.WriteFile(filepath.Join(dir, name), []byte(content), 0644))
	}
	for _, args := range [][]string{{"init", "--quiet"}, {"add", "-A"}, {"commit", "--quiet", "-m", "x"}} {
		cmd := exec.Command("git", args...)
		cmd.Dir = dir
		cmd.Env = append(os.Environ(), "GIT_CONFIG_NOSYSTEM=1", "GIT_CONFIG_GLOBAL=/dev/null",
			"GIT_AUTHOR_NAME=t", "GIT_AUTHOR_EMAIL=t@example.com", "GIT_COMMITTER_NAME=t", "GIT_COMMITTER_EMAIL=t@example.com")
		out, err := cmd.CombinedOutput()
		assert.NoError(t, err, string(out))
	}
}

func TestWasmPackRoutes(t *testing.T) {
	base := t.TempDir()
	createRepo(t, filepath.Join(base, "minetest_game"), map[string]string{"game.conf": "name = x\n"})
	createRepo(t, filepath.Join(base, "promise"), map[string]string{"mod.conf": "name = promise\n"})

	old_protocols, old_builtin := core.WasmPackGitProtocols, core.WasmPackBuiltinMods
	core.WasmPackGitProtocols = "file"
	core.WasmPackBuiltinMods = []string{"file://" + filepath.Join(base, "promise")}
	t.Cleanup(func() { core.WasmPackGitProtocols, core.WasmPackBuiltinMods = old_protocols, old_builtin })

	cfg := types.CreateConfig()
	cfg.WasmPackCacheDir = t.TempDir()
	a, _, err := api.NewApi(testutils.CreateTestDatabase(t), cfg)
	assert.NoError(t, err)

	owner := testutils.CreateUser(a.UserRepo, t, nil)
	other := testutils.CreateUser(a.UserRepo, t, nil)
	schema := testutils.CreateSchema(a.SchemaRepo, t, owner, nil)

	// one token per user: the config pack (and its etag) contains the token, a new jwt per request
	// would change the etag whenever the expiry (second resolution) differs
	tokens := map[string]string{}
	call := func(handler api.SecureHandler, user *types.User, if_none_match string) *httptest.ResponseRecorder {
		r := mux.SetURLVars(httptest.NewRequest("GET", "http://", nil), map[string]string{"schema_uid": schema.UID})
		if if_none_match != "" {
			r.Header.Set("If-None-Match", if_none_match)
		}
		w := httptest.NewRecorder()
		if tokens[user.UID] == "" {
			loginWith(t, r, user)
			tokens[user.UID] = r.Header.Get("Authorization")
		}
		r.Header.Set("Authorization", tokens[user.UID])
		a.Secure(handler)(w, r)
		return w
	}

	// no dependencies yet
	assert.Equal(t, 404, call(a.GetWasmPack, owner, "").Code)
	assert.Equal(t, 404, call(a.GetWasmConfigPack, owner, "").Code)
	assert.Equal(t, 404, call(a.GetWasmBuiltinPack, owner, "").Code)
	// the dependency url is validated by the api, local repositories are only possible directly in the database
	assert.NoError(t, a.SchemaModDependencyRepo.Create(&types.SchemaModDependency{
		SchemaUID: schema.UID, ModType: types.SchemaModDependencyTypeGame, GitURL: "file://" + filepath.Join(base, "minetest_game"),
	}))

	// owner only
	assert.Equal(t, 403, call(a.GetWasmPack, other, "").Code)
	assert.Equal(t, 403, call(a.GetWasmConfigPack, other, "").Code)
	assert.Equal(t, 403, call(a.GetWasmBuiltinPack, other, "").Code)

	for _, handler := range []api.SecureHandler{a.GetWasmPack, a.GetWasmBuiltinPack, a.GetWasmConfigPack} {
		w := call(handler, owner, "")
		assert.Equal(t, 200, w.Code)
		assert.Equal(t, "application/zstd", w.Header().Get("Content-Type"))
		assert.Equal(t, "private, no-cache", w.Header().Get("Cache-Control"))
		etag := w.Header().Get("ETag")
		assert.NotEmpty(t, etag)
		assert.Greater(t, w.Body.Len(), 0)

		// unchanged: not modified
		w = call(handler, owner, etag)
		assert.Equal(t, 304, w.Code)
		assert.Equal(t, 0, w.Body.Len())
	}
}
