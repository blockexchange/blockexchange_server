package core_test

import (
	"blockexchange/core"
	"blockexchange/metrics"
	"blockexchange/types"
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func gitCacheRepos(t *testing.T, cache_dir string) []string {
	entries, err := os.ReadDir(filepath.Join(cache_dir, "git"))
	if os.IsNotExist(err) {
		return nil
	}
	assert.NoError(t, err)
	names := []string{}
	for _, e := range entries {
		names = append(names, e.Name())
	}
	return names
}

func TestWasmGitCacheReuse(t *testing.T) {
	cache_dir := t.TempDir()
	c := newPackCoreDir(t, cache_dir)
	base := t.TempDir()
	a := filepath.Join(base, "a")
	commit_a := commitFiles(t, a, map[string]string{"init.lua": "-- a\n"})
	b := filepath.Join(base, "b")
	commit_b := commitFiles(t, b, map[string]string{"init.lua": "-- b\n"})

	remote := delta(metrics.WasmGitFetches.WithLabelValues("remote"))
	cached := delta(metrics.WasmGitFetches.WithLabelValues("cached"))

	_, err := c.GetWasmPack(context.Background(), core.WasmPackKindDependencies, []*types.SchemaModDependency{
		dep(types.SchemaModDependencyTypeMod, a, commit_a),
	})
	assert.NoError(t, err)
	assert.Equal(t, 1.0, remote())
	assert.Len(t, gitCacheRepos(t, cache_dir), 1)

	// another pack with the same commit of "a": read from the local repository, the remote isn't needed anymore
	assert.NoError(t, os.RemoveAll(a))
	p, err := c.GetWasmPack(context.Background(), core.WasmPackKindDependencies, []*types.SchemaModDependency{
		dep(types.SchemaModDependencyTypeMod, a, commit_a),
		dep(types.SchemaModDependencyTypeMod, b, commit_b),
	})
	assert.NoError(t, err)
	assert.Equal(t, 2.0, remote())
	assert.Equal(t, 1.0, cached())
	entries := readPack(t, p)
	assert.Equal(t, "-- a\n", entries["luanti/worlds/world/worldmods/a/init.lua"])
	assert.Equal(t, "-- b\n", entries["luanti/worlds/world/worldmods/b/init.lua"])
	assert.Len(t, gitCacheRepos(t, cache_dir), 2)

	// a new commit on "b": only that commit is fetched into the existing repository
	core.WasmPackResolveTTL = 0
	commitFiles(t, b, map[string]string{"init.lua": "-- b2\n"})
	p, err = c.GetWasmPack(context.Background(), core.WasmPackKindDependencies, []*types.SchemaModDependency{
		dep(types.SchemaModDependencyTypeMod, b, "latest"),
	})
	assert.NoError(t, err)
	assert.Equal(t, 3.0, remote())
	assert.Equal(t, "-- b2\n", readPack(t, p)["luanti/worlds/world/worldmods/b/init.lua"])
	assert.Len(t, gitCacheRepos(t, cache_dir), 2)
}

func TestWasmGitCacheEviction(t *testing.T) {
	cache_dir := t.TempDir()
	c := newPackCoreDir(t, cache_dir)
	old := core.WasmGitCacheMaxSize
	core.WasmGitCacheMaxSize = 1
	t.Cleanup(func() { core.WasmGitCacheMaxSize = old })

	a := filepath.Join(t.TempDir(), "a")
	commitFiles(t, a, map[string]string{"init.lua": "-- a\n"})
	deps := []*types.SchemaModDependency{dep(types.SchemaModDependencyTypeMod, a, "latest")}

	evictions := delta(metrics.WasmGitEvictions)
	_, err := c.GetWasmPack(context.Background(), core.WasmPackKindDependencies, deps)
	assert.NoError(t, err)
	assert.Equal(t, 1.0, evictions())
	assert.Len(t, gitCacheRepos(t, cache_dir), 0)

	// fetched again on the next build
	core.WasmPackResolveTTL = 0
	commitFiles(t, a, map[string]string{"init.lua": "-- a2\n"})
	p, err := c.GetWasmPack(context.Background(), core.WasmPackKindDependencies, deps)
	assert.NoError(t, err)
	assert.Equal(t, "-- a2\n", readPack(t, p)["luanti/worlds/world/worldmods/a/init.lua"])
}

func TestWasmGitArchiveAttributes(t *testing.T) {
	c := newPackCore(t)
	a := filepath.Join(t.TempDir(), "a")
	commitFiles(t, a, map[string]string{
		".gitattributes": "/spec export-ignore\n",
		"init.lua":       "-- a\n",
		"spec/test.lua":  "-- test\n",
		"tool.sh":        "#!/bin/sh\n",
	})
	assert.NoError(t, os.Chmod(filepath.Join(a, "tool.sh"), 0755))
	git(t, a, "add", "-A")
	git(t, a, "commit", "--quiet", "-m", "exec")

	p, err := c.GetWasmPack(context.Background(), core.WasmPackKindDependencies, []*types.SchemaModDependency{
		dep(types.SchemaModDependencyTypeMod, a, "latest"),
	})
	assert.NoError(t, err)
	entries := readPack(t, p)
	assert.Contains(t, entries, "luanti/worlds/world/worldmods/a/init.lua")
	assert.Contains(t, entries, "luanti/worlds/world/worldmods/a/tool.sh")
	for name := range entries {
		assert.NotContains(t, name, "spec", name)
	}

	modes := readPackModes(t, p)
	assert.Equal(t, int64(0755), modes["luanti/worlds/world/worldmods/a/tool.sh"])
	assert.Equal(t, int64(0644), modes["luanti/worlds/world/worldmods/a/init.lua"])
	assert.Equal(t, int64(0755), modes["luanti/worlds/world/worldmods/a/"])
}
