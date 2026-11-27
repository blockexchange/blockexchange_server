package core_test

import (
	"archive/tar"
	"blockexchange/core"
	"blockexchange/types"
	"bytes"
	"context"
	"encoding/json"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

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

// commits the files to the (new) local repository, returns the commit hash
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

// archive entries: name -> content ("" for directories)
func readPackData(t *testing.T, data []byte) map[string]string {
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

func readPack(t *testing.T, pack *core.WasmPack) map[string]string {
	data, err := os.ReadFile(pack.Path)
	assert.NoError(t, err)
	return readPackData(t, data)
}

// core with its own pack cache, local repositories allowed and the given builtin mods (no network access in tests)
func newPackCore(t *testing.T, builtin_mods ...string) *core.Core {
	return newPackCoreDir(t, t.TempDir(), builtin_mods...)
}

// core with the pack cache in cache_dir (git repositories in cache_dir/git)
func newPackCoreDir(t *testing.T, cache_dir string, builtin_mods ...string) *core.Core {
	old_protocols, old_builtin, old_ttl := core.WasmPackGitProtocols, core.WasmPackBuiltinMods, core.WasmPackResolveTTL
	core.WasmPackGitProtocols = "file"
	core.WasmPackBuiltinMods = builtin_mods
	t.Cleanup(func() {
		core.WasmPackGitProtocols, core.WasmPackBuiltinMods, core.WasmPackResolveTTL = old_protocols, old_builtin, old_ttl
	})

	cfg := types.CreateConfig()
	cfg.BaseURL = "https://bx.example.com"
	cfg.WasmPackCacheDir = cache_dir
	return core.New(cfg, nil)
}

func dep(mod_type types.SchemaModDependencyType, dir, version string) *types.SchemaModDependency {
	return &types.SchemaModDependency{ModType: mod_type, GitURL: "file://" + dir, Version: version}
}

func TestWasmPackDependencies(t *testing.T) {
	c := newPackCore(t)
	base := t.TempDir()

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

	pack, err := c.GetWasmPack(context.Background(), core.WasmPackKindDependencies, []*types.SchemaModDependency{
		dep(types.SchemaModDependencyTypeGame, game_dir, "latest"),
		dep(types.SchemaModDependencyTypeMod, mod_dir, first),
		dep(types.SchemaModDependencyTypeMod, other_dir, "dev"),
	})
	assert.NoError(t, err)
	assert.Equal(t, "my_game", pack.Manifest.GameID)
	assert.ElementsMatch(t, []string{"some_mod", "other_mod"}, pack.Manifest.Mods)
	assert.True(t, strings.HasPrefix(pack.ETag, `"dependencies-`))

	entries := readPack(t, pack)
	assert.Contains(t, entries, "luanti/worlds/world/worldmods/")
	assert.Equal(t, "title = My game\n", entries["luanti/worlds/world/game/game.conf"])
	assert.Equal(t, "-- x\n", entries["luanti/worlds/world/game/mods/x/init.lua"])
	assert.Equal(t, "-- v1\n", entries["luanti/worlds/world/worldmods/some_mod/init.lua"])
	assert.Equal(t, "-- dev\n", entries["luanti/worlds/world/worldmods/other_mod/init.lua"])
	assert.NotContains(t, entries, "luanti/minetest.conf")
	for name := range entries {
		assert.NotContains(t, name, ".git/", name)
		assert.NotContains(t, name, "passwd", name)
	}
}

func TestWasmPackCache(t *testing.T) {
	c := newPackCore(t)
	base := t.TempDir()
	game_dir := filepath.Join(base, "game")
	commitFiles(t, game_dir, map[string]string{"game.conf": "v1\n"})
	mod_dir := filepath.Join(base, "mod")
	mod_commit := commitFiles(t, mod_dir, map[string]string{"init.lua": "-- mod\n"})

	deps := []*types.SchemaModDependency{
		dep(types.SchemaModDependencyTypeGame, game_dir, "latest"),
		dep(types.SchemaModDependencyTypeMod, mod_dir, mod_commit),
	}
	p1, err := c.GetWasmPack(context.Background(), core.WasmPackKindDependencies, deps)
	assert.NoError(t, err)

	// cache hit: same pack, also in another order and without the mod repository
	assert.NoError(t, os.RemoveAll(mod_dir))
	p2, err := c.GetWasmPack(context.Background(), core.WasmPackKindDependencies, []*types.SchemaModDependency{deps[1], deps[0]})
	assert.NoError(t, err)
	assert.Equal(t, p1.Path, p2.Path)
	assert.Equal(t, p1.ETag, p2.ETag)

	// new commit on "latest": reused until the resolve ttl expires
	game_v2 := commitFiles(t, game_dir, map[string]string{"game.conf": "v2\n"})
	p3, err := c.GetWasmPack(context.Background(), core.WasmPackKindDependencies, deps[:1])
	assert.NoError(t, err)
	assert.Equal(t, "v1\n", readPack(t, p3)["luanti/worlds/world/game/game.conf"])

	core.WasmPackResolveTTL = 0
	p4, err := c.GetWasmPack(context.Background(), core.WasmPackKindDependencies, deps[:1])
	assert.NoError(t, err)
	assert.NotEqual(t, p3.ETag, p4.ETag)
	assert.Equal(t, "v2\n", readPack(t, p4)["luanti/worlds/world/game/game.conf"])

	// pinned to the new commit: same content, same pack
	p5, err := c.GetWasmPack(context.Background(), core.WasmPackKindDependencies, []*types.SchemaModDependency{
		dep(types.SchemaModDependencyTypeGame, game_dir, game_v2),
	})
	assert.NoError(t, err)
	assert.Equal(t, p4.ETag, p5.ETag)
}

func TestWasmPackConcurrent(t *testing.T) {
	c := newPackCore(t)
	game_dir := filepath.Join(t.TempDir(), "game")
	commitFiles(t, game_dir, map[string]string{"game.conf": "x\n"})
	deps := []*types.SchemaModDependency{dep(types.SchemaModDependencyTypeGame, game_dir, "latest")}

	wg := sync.WaitGroup{}
	etags := make([]string, 8)
	for i := range etags {
		wg.Go(func() {
			p, err := c.GetWasmPack(context.Background(), core.WasmPackKindDependencies, deps)
			assert.NoError(t, err)
			if p != nil {
				etags[i] = p.ETag
			}
		})
	}
	wg.Wait()
	for _, etag := range etags {
		assert.Equal(t, etags[0], etag)
	}

	// an already cancelled request fails
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	core.WasmPackResolveTTL = 0
	_, err := c.GetWasmPack(ctx, core.WasmPackKindDependencies, deps)
	assert.Error(t, err)
}

func TestWasmPackEviction(t *testing.T) {
	c := newPackCore(t)
	old := core.WasmPackCacheMaxSize
	core.WasmPackCacheMaxSize = 1
	t.Cleanup(func() { core.WasmPackCacheMaxSize = old })

	base := t.TempDir()
	a := filepath.Join(base, "a")
	commitFiles(t, a, map[string]string{"init.lua": "-- a\n"})
	b := filepath.Join(base, "b")
	commitFiles(t, b, map[string]string{"init.lua": "-- b\n"})

	pa, err := c.GetWasmPack(context.Background(), core.WasmPackKindDependencies, []*types.SchemaModDependency{dep(types.SchemaModDependencyTypeMod, a, "latest")})
	assert.NoError(t, err)
	time.Sleep(10 * time.Millisecond)
	pb, err := c.GetWasmPack(context.Background(), core.WasmPackKindDependencies, []*types.SchemaModDependency{dep(types.SchemaModDependencyTypeMod, b, "latest")})
	assert.NoError(t, err)

	// over the limit: the older pack is removed, the new one kept
	_, err = os.Stat(pa.Path)
	assert.True(t, os.IsNotExist(err))
	_, err = os.Stat(pb.Path)
	assert.NoError(t, err)
}

func TestWasmPackTagVersion(t *testing.T) {
	c := newPackCore(t)
	mod_dir := filepath.Join(t.TempDir(), "mod")
	commitFiles(t, mod_dir, map[string]string{"init.lua": "-- tagged\n"})
	git(t, mod_dir, "tag", "-a", "v1.0", "-m", "release")
	commitFiles(t, mod_dir, map[string]string{"init.lua": "-- newer\n"})

	p, err := c.GetWasmPack(context.Background(), core.WasmPackKindDependencies, []*types.SchemaModDependency{dep(types.SchemaModDependencyTypeMod, mod_dir, "v1.0")})
	assert.NoError(t, err)
	assert.Equal(t, "-- tagged\n", readPack(t, p)["luanti/worlds/world/worldmods/mod/init.lua"])
}

func TestWasmPackInvalid(t *testing.T) {
	c := newPackCore(t)
	base := t.TempDir()
	a := filepath.Join(base, "a")
	commitFiles(t, a, map[string]string{"mod.conf": "name = same\n"})
	b := filepath.Join(base, "b")
	commitFiles(t, b, map[string]string{"mod.conf": "name = same\n"})

	get := func(deps ...*types.SchemaModDependency) error {
		_, err := c.GetWasmPack(context.Background(), core.WasmPackKindDependencies, deps)
		return err
	}

	// two games
	assert.ErrorIs(t, get(dep(types.SchemaModDependencyTypeGame, a, "latest"), dep(types.SchemaModDependencyTypeGame, b, "latest")), core.ErrWasmPackInvalid)
	// same mod folder
	assert.ErrorIs(t, get(dep(types.SchemaModDependencyTypeMod, a, "latest"), dep(types.SchemaModDependencyTypeMod, b, "latest")), core.ErrWasmPackInvalid)
	// option injection via version
	assert.ErrorIs(t, get(dep(types.SchemaModDependencyTypeMod, a, "--upload-pack=touch /tmp/x")), core.ErrWasmPackInvalid)
	// unknown version
	assert.ErrorIs(t, get(dep(types.SchemaModDependencyTypeMod, a, "does-not-exist")), core.ErrWasmPackInvalid)

	// file:// is not allowed by default
	core.WasmPackGitProtocols = "http:https"
	core.WasmPackResolveTTL = 0
	err := get(dep(types.SchemaModDependencyTypeMod, b, "latest"))
	assert.Error(t, err)
	assert.NotErrorIs(t, err, core.ErrWasmPackInvalid)
}

func TestWasmBuiltinAndConfigPack(t *testing.T) {
	base := t.TempDir()
	builtin_bx := filepath.Join(base, "builtin", "blockexchange")
	commitFiles(t, builtin_bx, map[string]string{"mod.conf": "name = blockexchange\n", "init.lua": "-- builtin\n"})
	builtin_promise := filepath.Join(base, "builtin", "promise")
	commitFiles(t, builtin_promise, map[string]string{"mod.conf": "name = promise\n", "init.lua": "-- promise\n"})
	c := newPackCore(t, "file://"+builtin_bx, "file://"+builtin_promise)

	// all builtin mods in worldmods
	full, err := c.GetWasmBuiltinPack(context.Background(), []string{"other_mod"})
	assert.NoError(t, err)
	assert.True(t, strings.HasPrefix(full.ETag, `"builtin-`))
	assert.ElementsMatch(t, []string{"blockexchange", "promise"}, full.Manifest.Mods)
	entries := readPack(t, full)
	assert.Equal(t, "-- builtin\n", entries["luanti/worlds/world/worldmods/blockexchange/init.lua"])
	assert.Equal(t, "-- promise\n", entries["luanti/worlds/world/worldmods/promise/init.lua"])
	assert.NotContains(t, entries, "luanti/mods/")

	// no overlap: the full pack is shared
	full2, err := c.GetWasmBuiltinPack(context.Background(), nil)
	assert.NoError(t, err)
	assert.Equal(t, full.ETag, full2.ETag)

	// the schematic brings its own blockexchange mod: left out of the builtin pack
	game_dir := filepath.Join(base, "minetest_game")
	commitFiles(t, game_dir, map[string]string{"game.conf": "x\n"})
	fork := filepath.Join(base, "blockexchange_fork")
	commitFiles(t, fork, map[string]string{"mod.conf": "name = blockexchange\n", "init.lua": "-- fork\n"})
	deps, err := c.GetWasmPack(context.Background(), core.WasmPackKindDependencies, []*types.SchemaModDependency{
		dep(types.SchemaModDependencyTypeGame, game_dir, "latest"),
		dep(types.SchemaModDependencyTypeMod, fork, "latest"),
	})
	assert.NoError(t, err)
	assert.Equal(t, "-- fork\n", readPack(t, deps)["luanti/worlds/world/worldmods/blockexchange/init.lua"])

	reduced, err := c.GetWasmBuiltinPack(context.Background(), deps.Manifest.Mods)
	assert.NoError(t, err)
	assert.NotEqual(t, full.ETag, reduced.ETag)
	assert.Equal(t, []string{"promise"}, reduced.Manifest.Mods)
	entries = readPack(t, reduced)
	assert.Contains(t, entries, "luanti/worlds/world/worldmods/promise/init.lua")
	for name := range entries {
		assert.NotContains(t, name, "worldmods/blockexchange", name)
	}

	data, err := c.WasmConfigPack("a681b2e3-f44a-47e5-b862-5c4714709040", "my.jwt.token", deps.Manifest)
	assert.NoError(t, err)
	entries = readPackData(t, data)

	// the embedded files (core/wasm_config) with the generated values
	embedded := func(name string) string {
		data, err := os.ReadFile(filepath.Join("wasm_config", name))
		assert.NoError(t, err)
		return string(data)
	}

	world_mt := entries["luanti/worlds/world/world.mt"]
	assert.Contains(t, world_mt, "gameid = minetest_game\n")
	// no world_name: the folder name "world" is the world name (--worldname world)
	assert.NotContains(t, world_mt, "world_name")
	assert.Contains(t, world_mt, embedded("world.mt"))
	// worldmods are always loaded
	assert.NotContains(t, world_mt, "load_mod_")

	map_meta, found := entries["luanti/worlds/world/map_meta.txt"]
	assert.True(t, found)
	assert.Equal(t, embedded("map_meta.txt"), map_meta)

	provision := &core.WasmProvision{}
	assert.NoError(t, json.Unmarshal([]byte(entries["luanti/worlds/world/bx_provision.json"]), provision))
	assert.Equal(t, "a681b2e3-f44a-47e5-b862-5c4714709040", provision.SchemaUID)
	assert.Equal(t, core.WasmProvisionPosition{}, provision.Position)
	assert.Equal(t, "my.jwt.token", provision.Token)

	conf := entries["luanti/minetest.conf"]
	assert.Contains(t, conf, embedded("minetest.conf"))
	assert.Contains(t, conf, "blockexchange.url = https://bx.example.com\n")
}

// archive entries: name -> file mode
func readPackModes(t *testing.T, pack *core.WasmPack) map[string]int64 {
	f, err := os.Open(pack.Path)
	assert.NoError(t, err)
	defer f.Close()
	zr, err := zstd.NewReader(f)
	assert.NoError(t, err)
	defer zr.Close()

	modes := map[string]int64{}
	tr := tar.NewReader(zr)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		assert.NoError(t, err)
		modes[hdr.Name] = hdr.Mode
	}
	return modes
}
