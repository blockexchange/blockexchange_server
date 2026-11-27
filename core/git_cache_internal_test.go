package core

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGitPruneCommits(t *testing.T) {
	old := WasmPackGitProtocols
	WasmPackGitProtocols = "file"
	t.Cleanup(func() { WasmPackGitProtocols = old })

	src := t.TempDir()
	gitT := func(env []string, args ...string) string {
		cmd := exec.Command("git", args...)
		cmd.Dir = src
		cmd.Env = append(os.Environ(),
			"GIT_CONFIG_NOSYSTEM=1", "GIT_CONFIG_GLOBAL=/dev/null",
			"GIT_AUTHOR_NAME=test", "GIT_AUTHOR_EMAIL=test@example.com",
			"GIT_COMMITTER_NAME=test", "GIT_COMMITTER_EMAIL=test@example.com",
		)
		cmd.Env = append(cmd.Env, env...)
		out, err := cmd.CombinedOutput()
		assert.NoError(t, err, string(out))
		return strings.TrimSpace(string(out))
	}
	gitT(nil, "init", "--quiet")

	ctx := context.Background()
	cache := newGitRepoCache(t.TempDir())
	url := "file://" + src
	repo := cache.acquire(url)

	commits := []string{}
	for i := range 11 {
		assert.NoError(t, os.WriteFile(filepath.Join(src, "init.lua"), []byte(strings.Repeat("x", i+1)), 0644))
		gitT(nil, "add", "-A")
		// distinct committer dates for the ordering
		gitT([]string{fmt.Sprintf("GIT_COMMITTER_DATE=2020-01-01T00:00:%02d", i)}, "commit", "--quiet", "-m", "c")
		commit := gitT(nil, "rev-parse", "HEAD")
		commits = append(commits, commit)
		assert.NoError(t, cache.ensureCommit(ctx, repo, url, commit))
	}
	cache.release(repo)

	// in use: not pruned
	cache.acquire(url)
	cache.maintain(ctx)
	refs, err := runGit(ctx, repo.dir, "for-each-ref", "--format=%(objectname)", gitCacheRefs)
	assert.NoError(t, err)
	assert.ElementsMatch(t, commits, strings.Fields(refs))
	cache.release(repo)

	old_max := WasmGitCacheMaxCommits
	WasmGitCacheMaxCommits = 5
	t.Cleanup(func() { WasmGitCacheMaxCommits = old_max })
	cache.maintain(ctx)

	refs, err = runGit(ctx, repo.dir, "for-each-ref", "--format=%(objectname)", gitCacheRefs)
	assert.NoError(t, err)
	assert.ElementsMatch(t, commits[6:], strings.Fields(refs))

	// the old commits are gone, the newest can still be archived
	_, err = runGit(ctx, repo.dir, "cat-file", "-e", commits[0]+"^{commit}")
	assert.Error(t, err)
	_, err = runGit(ctx, repo.dir, "archive", "--format=tar", commits[10])
	assert.NoError(t, err)
}
