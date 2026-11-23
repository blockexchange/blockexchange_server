package core_test

import (
	"archive/tar"
	"blockexchange/core"
	"blockexchange/types"
	"bytes"
	"context"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/klauspost/compress/zstd"
	"github.com/stretchr/testify/assert"
)

func git(t *testing.T, dir string, args ...string) string {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(),
		"GIT_CONFIG_NOSYSTEM=1", "GIT_CONFIG_GLOBAL=/dev/null",
		"GIT_AUTHOR_NAME=test", "GIT_AUTHOR_EMAIL=test@example.com",
		"GIT_COMMITTER_NAME=test", "GIT_COMMITTER_EMAIL=test@example.com",
	)
	out, err := cmd.CombinedOutput()
	assert.NoError(t, err, string(out))
	return strings.TrimSpace(string(out))
}

// creates a local repository with the given files committed, returns the commit hash
func commitFiles(t *testing.T, dir string, files map[string]string) string {
	if _, err := os.Stat(filepath.Join(dir, ".git")); err != nil {
		assert.NoError(t, os.MkdirAll(dir, 0755))
		git(t, dir, "init", "--quiet", "--initial-branch=main")
	}
	for name, content := range files {
		p := filepath.Join(dir, name)
		assert.NoError(t, os.MkdirAll(filepath.Dir(p), 0755))
		assert.NoError(t, os.WriteFile(p, []byte(content), 0644))
	}
	git(t, dir, "add", "-A")
	git(t, dir, "commit", "--quiet", "-m", "commit")
	return git(t, dir, "rev-parse", "HEAD")
}

// returns the archive entries: name -> content ("" for directories)
func readPack(t *testing.T, data []byte) map[string]string {
	zr, err := zstd.NewReader(bytes.NewReader(data))
	assert.NoError(t, err)
	defer zr.Close()

	entries := map[string]string{}
	tr := tar.NewReader(zr)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		assert.NoError(t, err)
		content, err := io.ReadAll(tr)
		assert.NoError(t, err)
		entries[hdr.Name] = string(content)
	}
	return entries
}

// allows local repositories and replaces the builtin mods (no network access in tests)
func allowFileProtocol(t *testing.T, builtin_mods ...string) {
	old := core.WasmPackGitProtocols
	old_builtin := core.WasmPackBuiltinMods
	core.WasmPackGitProtocols = "file"
	core.WasmPackBuiltinMods = builtin_mods
	t.Cleanup(func() {
		core.WasmPackGitProtocols = old
		core.WasmPackBuiltinMods = old_builtin
	})
}

func TestWasmPack(t *testing.T) {
	allowFileProtocol(t)
	base := t.TempDir()
	c := core.New(types.CreateConfig(), nil)

	game_dir := filepath.Join(base, "my_game.git")
	commitFiles(t, game_dir, map[string]string{"game.conf": "title = My game\n", "mods/x/init.lua": "-- x\n"})

	// mod with a name in mod.conf that differs from the repository name, pinned to the first commit
	mod_dir := filepath.Join(base, "Some-Mod")
	first := commitFiles(t, mod_dir, map[string]string{"mod.conf": "name = some_mod\n", "init.lua": "-- v1\n"})
	commitFiles(t, mod_dir, map[string]string{"init.lua": "-- v2\n"})
	assert.NoError(t, os.Symlink("/etc/passwd", filepath.Join(mod_dir, "passwd")))
	git(t, mod_dir, "add", "-A")
	git(t, mod_dir, "commit", "--quiet", "-m", "symlink")

	// mod without mod.conf, pinned to a branch
	other_dir := filepath.Join(base, "other-mod.git")
	commitFiles(t, other_dir, map[string]string{"init.lua": "-- main\n"})
	git(t, other_dir, "checkout", "--quiet", "-b", "dev")
	commitFiles(t, other_dir, map[string]string{"init.lua": "-- dev\n"})
	git(t, other_dir, "checkout", "--quiet", "main")

	deps := []*types.SchemaModDependency{
		{ModType: types.SchemaModDependencyTypeGame, GitURL: "file://" + game_dir, Version: "latest"},
		{ModType: types.SchemaModDependencyTypeMod, GitURL: "file://" + mod_dir, Version: first},
		{ModType: types.SchemaModDependencyTypeMod, GitURL: "file://" + other_dir, Version: "dev"},
	}

	pack, err := c.PrepareWasmPack(context.Background(), "my_schema", deps)
	assert.NoError(t, err)
	buf := bytes.NewBuffer(nil)
	assert.NoError(t, pack.Write(buf))
	assert.NoError(t, pack.Close())

	entries := readPack(t, buf.Bytes())
	assert.Contains(t, entries, "luanti/worlds/world/")
	assert.Contains(t, entries, "luanti/worlds/world/worldmods/")
	assert.Equal(t, "title = My game\n", entries["luanti/worlds/world/game/game.conf"])
	assert.Equal(t, "-- x\n", entries["luanti/worlds/world/game/mods/x/init.lua"])
	assert.Equal(t, "-- v1\n", entries["luanti/worlds/world/worldmods/some_mod/init.lua"])
	assert.Equal(t, "-- dev\n", entries["luanti/worlds/world/worldmods/other_mod/init.lua"])

	assert.Contains(t, entries["luanti/minetest.conf"], "creative_mode = true\n")
	world_mt := entries["luanti/worlds/world/world.mt"]
	assert.Contains(t, world_mt, "gameid = my_game\n")
	assert.Contains(t, world_mt, "world_name = my_schema\n")
	assert.Contains(t, world_mt, "backend = sqlite3\n")

	for name := range entries {
		assert.NotContains(t, name, ".git/", name)
		assert.NotContains(t, name, "passwd", name)
	}
}

func TestWasmPackInvalid(t *testing.T) {
	allowFileProtocol(t)
	base := t.TempDir()
	c := core.New(types.CreateConfig(), nil)

	a := filepath.Join(base, "a")
	commitFiles(t, a, map[string]string{"mod.conf": "name = same\n"})
	b := filepath.Join(base, "b")
	commitFiles(t, b, map[string]string{"mod.conf": "name = same\n"})

	// two games
	_, err := c.PrepareWasmPack(context.Background(), "x", []*types.SchemaModDependency{
		{ModType: types.SchemaModDependencyTypeGame, GitURL: "file://" + a},
		{ModType: types.SchemaModDependencyTypeGame, GitURL: "file://" + b},
	})
	assert.ErrorIs(t, err, core.ErrWasmPackInvalid)

	// same mod folder
	_, err = c.PrepareWasmPack(context.Background(), "x", []*types.SchemaModDependency{
		{ModType: types.SchemaModDependencyTypeMod, GitURL: "file://" + a, Version: "latest"},
		{ModType: types.SchemaModDependencyTypeMod, GitURL: "file://" + b, Version: "latest"},
	})
	assert.ErrorIs(t, err, core.ErrWasmPackInvalid)

	// option injection via version
	_, err = c.PrepareWasmPack(context.Background(), "x", []*types.SchemaModDependency{
		{ModType: types.SchemaModDependencyTypeMod, GitURL: "file://" + a, Version: "--upload-pack=touch /tmp/x"},
	})
	assert.ErrorIs(t, err, core.ErrWasmPackInvalid)

	// unknown version
	_, err = c.PrepareWasmPack(context.Background(), "x", []*types.SchemaModDependency{
		{ModType: types.SchemaModDependencyTypeMod, GitURL: "file://" + a, Version: "does-not-exist"},
	})
	assert.Error(t, err)
	assert.NotErrorIs(t, err, core.ErrWasmPackInvalid)

	// file:// is not allowed by default
	core.WasmPackGitProtocols = "http:https"
	_, err = c.PrepareWasmPack(context.Background(), "x", []*types.SchemaModDependency{
		{ModType: types.SchemaModDependencyTypeMod, GitURL: "file://" + a, Version: "latest"},
	})
	assert.Error(t, err)
}

func TestWasmPackBuiltinMods(t *testing.T) {
	base := t.TempDir()
	cfg := types.CreateConfig()
	cfg.BaseURL = "https://bx.example.com"
	c := core.New(cfg, nil)

	builtin_bx := filepath.Join(base, "builtin", "blockexchange")
	commitFiles(t, builtin_bx, map[string]string{"mod.conf": "name = blockexchange\n", "init.lua": "-- builtin\n"})
	builtin_promise := filepath.Join(base, "builtin", "promise")
	commitFiles(t, builtin_promise, map[string]string{"mod.conf": "name = promise\n", "init.lua": "-- promise\n"})
	allowFileProtocol(t, "file://"+builtin_bx, "file://"+builtin_promise)

	// the schematic brings its own blockexchange mod, the builtin one is skipped
	own_bx := filepath.Join(base, "own", "blockexchange_fork")
	commitFiles(t, own_bx, map[string]string{"mod.conf": "name = blockexchange\n", "init.lua": "-- fork\n"})

	pack, err := c.PrepareWasmPack(context.Background(), "x", []*types.SchemaModDependency{
		{ModType: types.SchemaModDependencyTypeMod, GitURL: "file://" + own_bx, Version: "latest"},
	})
	assert.NoError(t, err)
	buf := bytes.NewBuffer(nil)
	assert.NoError(t, pack.Write(buf))
	assert.NoError(t, pack.Close())

	entries := readPack(t, buf.Bytes())
	assert.Equal(t, "-- fork\n", entries["luanti/worlds/world/worldmods/blockexchange/init.lua"])
	assert.Equal(t, "-- promise\n", entries["luanti/worlds/world/worldmods/promise/init.lua"])

	conf := entries["luanti/minetest.conf"]
	assert.Contains(t, conf, "secure.http_mods = blockexchange\n")
	assert.Contains(t, conf, "blockexchange.url = https://bx.example.com\n")
}
