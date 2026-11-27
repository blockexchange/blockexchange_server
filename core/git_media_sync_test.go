package core_test

import (
	"blockexchange/core"
	"blockexchange/testutils"
	"blockexchange/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGitMediaSync(t *testing.T) {
	old_protocols := core.WasmPackGitProtocols
	core.WasmPackGitProtocols = "file"
	t.Cleanup(func() { core.WasmPackGitProtocols = old_protocols })

	repos := testutils.CreateTestDatabase(t)
	cfg := types.CreateConfig()
	cfg.GitMediaCacheDir = t.TempDir()
	c := core.New(cfg, repos)

	dir := t.TempDir()
	commitFiles(t, dir, map[string]string{
		"init.lua":                   "",
		"textures/gitsync_a.png":     "a",
		"textures/sub/gitsync_b.png": "bb",
		"models/gitsync_c.obj":       "ccc",
		"sounds/gitsync_d.ogg":       "d",
	})

	// licenses are set, no contentdb lookup
	src := &types.GitMediaSource{GitURL: "file://" + dir, MediaLicense: "CC0-1.0", License: "MIT"}
	assert.NoError(t, repos.GitMediaSourceRepo.Create(src))
	t.Cleanup(func() { repos.GitMediaSourceRepo.Delete(src.ID) })

	assert.NoError(t, c.SyncGitMediaSource(src))
	assert.True(t, src.LastUpdatedSuccess, src.LastUpdatedMessage)
	names, err := repos.MediaRepo.GetMediafileNamesBySource(src.ID)
	assert.NoError(t, err)
	assert.Equal(t, []string{"gitsync_a.png", "gitsync_b.png", "gitsync_c.obj"}, names)
	mf, err := repos.MediaRepo.GetMediafileByName("gitsync_b.png")
	assert.NoError(t, err)
	assert.Equal(t, []byte("bb"), mf.Data)

	// update: changed and removed files
	git(t, dir, "rm", "--quiet", "models/gitsync_c.obj")
	commitFiles(t, dir, map[string]string{"textures/gitsync_a.png": "a2"})

	assert.NoError(t, c.SyncGitMediaSource(src))
	assert.True(t, src.LastUpdatedSuccess, src.LastUpdatedMessage)
	names, err = repos.MediaRepo.GetMediafileNamesBySource(src.ID)
	assert.NoError(t, err)
	assert.Equal(t, []string{"gitsync_a.png", "gitsync_b.png"}, names)
	mf, err = repos.MediaRepo.GetMediafileByName("gitsync_a.png")
	assert.NoError(t, err)
	assert.Equal(t, []byte("a2"), mf.Data)

	// failed fetch
	src.GitURL = "file:///nonexistent"
	assert.NoError(t, c.SyncGitMediaSource(src))
	assert.False(t, src.LastUpdatedSuccess)
}
