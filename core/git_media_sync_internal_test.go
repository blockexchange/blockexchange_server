package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsGitMediaPath(t *testing.T) {
	assert.True(t, isGitMediaPath("textures/default_stone.png"))
	assert.True(t, isGitMediaPath("textures/sub/default_stone.PNG"))
	assert.True(t, isGitMediaPath("models/torch.obj"))
	assert.True(t, isGitMediaPath("mymod/textures/x.png"))
	assert.True(t, isGitMediaPath("mods/mymod/models/x.b3d"))
	assert.False(t, isGitMediaPath("default_stone.png"))
	assert.False(t, isGitMediaPath("textures/readme.txt"))
	assert.False(t, isGitMediaPath("sounds/x.ogg"))
	assert.False(t, isGitMediaPath("screenshot.png"))
	assert.False(t, isGitMediaPath("mytextures/x.png"))
}

func TestParseGitMediaTree(t *testing.T) {
	out := "100644 blob aaa 10\ttextures/a.png\x00" +
		"100644 blob bbb 20\tmod2/textures/a.png\x00" + // duplicate name, the first one wins
		"100755 blob ccc 30\tmodels/b.obj\x00" +
		"120000 blob ddd 5\ttextures/link.png\x00" + // symlink
		"160000 commit eee -\ttextures/sub\x00" + // submodule
		"100644 blob fff 99999999999\ttextures/huge.png\x00" +
		"100644 blob ggg 1\tinit.lua\x00"
	entries, skipped := parseGitMediaTree(out)
	assert.Equal(t, 1, skipped)
	assert.Equal(t, []gitMediaEntry{
		{name: "a.png", object: "aaa", size: 10},
		{name: "b.obj", object: "ccc", size: 30},
	}, entries)
}

func TestFindCDBPackageByRepo(t *testing.T) {
	packages := []CDBVcsPackage{
		{Author: "a", Name: "nope", Repo: ""},
		{Author: "mt-mods", Name: "hangglider", Repo: "https://github.com/mt-mods/hangglider.git"},
	}
	assert.Equal(t, "mt-mods/hangglider", findCDBPackageByRepo(packages, "https://github.com/mt-mods/hangglider"))
	assert.Equal(t, "mt-mods/hangglider", findCDBPackageByRepo(packages, "http://GitHub.com/mt-mods/hangglider/"))
	assert.Equal(t, "", findCDBPackageByRepo(packages, "https://github.com/mt-mods/other"))
}
