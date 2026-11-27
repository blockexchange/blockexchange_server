package core

import (
	"archive/tar"
	"blockexchange/metrics"
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// local cache of the git repositories used in the wasm packs: one bare repository per url,
// only the needed commits are fetched (shallow) and the files are read from the object store without a checkout
//
//	<dir>/<sha256(url)>/        bare repository
//	<dir>/<sha256(url)>/bx-url  the url, for debugging
//	<dir>/<sha256(url)>/bx-used touched on every use, for the eviction

var (
	// total size of the cached repositories, the least recently used ones are removed first
	WasmGitCacheMaxSize = int64(50 * 1024 * 1024 * 1024)
	// commits kept per repository, the oldest ones are removed once there are twice as many
	WasmGitCacheMaxCommits = 5
)

const (
	gitCacheURLFile  = "bx-url"
	gitCacheUsedFile = "bx-used"
	gitCacheRefs     = "refs/bx/"
	// upper limit for mod.conf / modpack.conf
	gitConfMaxSize = 64 * 1024
)

type gitRepoCache struct {
	dir string

	// guards repos and the users/busy state of each repo
	mu    sync.Mutex
	cond  *sync.Cond
	repos map[string]*gitRepo
}

type gitRepo struct {
	dir string
	// serializes fetches into the repository
	fetch_lock sync.Mutex
	// number of builds using the repository, maintenance only runs without users
	users int
	// maintenance (pruning, removal) in progress, new users wait
	busy bool
}

func newGitRepoCache(dir string) *gitRepoCache {
	c := &gitRepoCache{dir: dir, repos: map[string]*gitRepo{}}
	c.cond = sync.NewCond(&c.mu)
	return c
}

func gitRepoKey(url string) string {
	h := sha256.Sum256([]byte(url))
	return hex.EncodeToString(h[:16])
}

// marks the repository as used, its commits stay available until release() is called
func (c *gitRepoCache) acquire(url string) *gitRepo {
	dir := filepath.Join(c.dir, gitRepoKey(url))
	c.mu.Lock()
	defer c.mu.Unlock()
	repo := c.repos[dir]
	if repo == nil {
		repo = &gitRepo{dir: dir}
		c.repos[dir] = repo
	}
	for repo.busy {
		c.cond.Wait()
	}
	repo.users++
	return repo
}

func (c *gitRepoCache) release(repo *gitRepo) {
	c.mu.Lock()
	repo.users--
	c.mu.Unlock()
	c.cond.Broadcast()
}

// claims the repository for maintenance, false if it is in use
func (c *gitRepoCache) tryMaintain(dir string) (*gitRepo, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	repo := c.repos[dir]
	if repo == nil {
		repo = &gitRepo{dir: dir}
		c.repos[dir] = repo
	}
	if repo.users > 0 || repo.busy {
		return nil, false
	}
	repo.busy = true
	return repo, true
}

func (c *gitRepoCache) doneMaintain(repo *gitRepo) {
	c.mu.Lock()
	repo.busy = false
	c.mu.Unlock()
	c.cond.Broadcast()
}

// makes sure the commit is in the local repository, fetches only that commit if it is missing
func (c *gitRepoCache) ensureCommit(ctx context.Context, repo *gitRepo, url, commit string) error {
	repo.fetch_lock.Lock()
	defer repo.fetch_lock.Unlock()

	created := false
	if _, err := os.Stat(filepath.Join(repo.dir, "HEAD")); err != nil {
		err = gitInitCacheRepo(ctx, repo.dir, url)
		if err != nil {
			os.RemoveAll(repo.dir)
			return err
		}
		created = true
	}
	now := time.Now()
	os.Chtimes(filepath.Join(repo.dir, gitCacheUsedFile), now, now)

	if !created {
		if _, err := runGit(ctx, repo.dir, "cat-file", "-e", commit+"^{commit}"); err == nil {
			metrics.WasmGitFetches.WithLabelValues("cached").Inc()
			return nil
		}
	}

	// the ref keeps the commit from being pruned, the existing refs are sent as "have"
	// so the remote only sends the objects that are not here yet
	_, err := runGit(ctx, repo.dir, "fetch", "--quiet", "--depth", "1", "--no-tags", "--no-write-fetch-head",
		"--end-of-options", url, "+"+commit+":"+gitCacheRefs+commit)
	if err != nil {
		metrics.WasmGitFetches.WithLabelValues("error").Inc()
		if created {
			os.RemoveAll(repo.dir)
		}
		return err
	}
	metrics.WasmGitFetches.WithLabelValues("remote").Inc()
	return nil
}

func gitInitCacheRepo(ctx context.Context, dir, url string) error {
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}
	_, err = runGit(ctx, dir, "init", "--bare", "--quiet")
	if err != nil {
		return err
	}
	for _, kv := range [][2]string{
		// garbage collection only in prune(), never in the background while a pack is read
		{"gc.auto", "0"},
		{"maintenance.auto", "false"},
		// fetched objects stay in packs: fewer files to walk for the size calculation
		{"fetch.unpackLimit", "1"},
	} {
		_, err = runGit(ctx, dir, "config", kv[0], kv[1])
		if err != nil {
			return err
		}
	}
	err = os.WriteFile(filepath.Join(dir, gitCacheURLFile), []byte(url+"\n"), 0644)
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, gitCacheUsedFile), nil, 0644)
}

// value of the key in mod.conf or modpack.conf (the first one with a valid name), empty if none
func gitReadConfValue(ctx context.Context, repo *gitRepo, commit, file, key string) string {
	out, err := runGit(ctx, repo.dir, "ls-tree", "-l", "-z", "--end-of-options", commit, "--", file)
	if err != nil {
		return ""
	}
	// "<mode> <type> <object> <size>\t<name>\x00"
	meta, name, found := strings.Cut(strings.TrimSuffix(out, "\x00"), "\t")
	fields := strings.Fields(meta)
	if !found || name != file || len(fields) != 4 || fields[1] != "blob" || !strings.HasPrefix(fields[0], "100") {
		// missing, symlink, submodule or directory
		return ""
	}
	size, err := strconv.ParseInt(fields[3], 10, 64)
	if err != nil || size > gitConfMaxSize {
		return ""
	}
	data, err := runGit(ctx, repo.dir, "cat-file", "blob", fields[2])
	if err != nil {
		return ""
	}
	return readConfValueFrom(strings.NewReader(data), key)
}

// copies the files of the commit into the archive below archive_dir ("git archive", applies export-ignore),
// skips symlinks and fails if more than budget bytes are written, returns the bytes written
func gitWriteTar(ctx context.Context, tw *tar.Writer, repo *gitRepo, commit, archive_dir string, budget int64) (int64, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	cmd := exec.CommandContext(ctx, "git", "archive", "--format=tar", "--end-of-options", commit)
	cmd.Dir = repo.dir
	cmd.Env = gitEnv()
	stderr := &bytes.Buffer{}
	cmd.Stderr = stderr
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return 0, err
	}
	err = cmd.Start()
	if err != nil {
		return 0, err
	}

	size, copy_err := copyGitArchive(tw, stdout, archive_dir, budget)
	if copy_err != nil {
		cancel()
		cmd.Wait()
		return 0, copy_err
	}
	// drain the trailing padding so git doesn't fail on a closed pipe
	io.Copy(io.Discard, stdout)
	err = cmd.Wait()
	if ctx.Err() != nil {
		return 0, ctx.Err()
	}
	if err != nil {
		return 0, fmt.Errorf("git archive failed: %s", strings.TrimSpace(stderr.String()))
	}
	return size, nil
}

func copyGitArchive(tw *tar.Writer, r io.Reader, archive_dir string, budget int64) (int64, error) {
	// the root folder isn't in the git archive
	err := tw.WriteHeader(&tar.Header{Typeflag: tar.TypeDir, Name: archive_dir + "/", Mode: 0755})
	if err != nil {
		return 0, err
	}

	tr := tar.NewReader(r)
	var size int64
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			return size, nil
		}
		if err != nil {
			return 0, err
		}

		name := path.Clean(hdr.Name)
		if name == "." || name == ".." || strings.HasPrefix(name, "../") || path.IsAbs(name) {
			continue
		}
		out := &tar.Header{Name: path.Join(archive_dir, name), ModTime: hdr.ModTime}
		switch hdr.Typeflag {
		case tar.TypeDir:
			out.Typeflag, out.Mode, out.Name = tar.TypeDir, 0755, out.Name+"/"
		case tar.TypeReg:
			out.Typeflag, out.Mode, out.Size = tar.TypeReg, 0644, hdr.Size
			if hdr.Mode&0100 != 0 {
				out.Mode = 0755
			}
			size += hdr.Size
			if size > budget {
				return 0, invalidPack("pack exceeds the size limit of %d bytes", WasmPackMaxSize)
			}
		default:
			// symlinks could point outside of the repository, pax global header with the commit id
			continue
		}

		err = tw.WriteHeader(out)
		if err != nil {
			return 0, err
		}
		if out.Typeflag == tar.TypeReg {
			_, err = io.Copy(tw, tr)
			if err != nil {
				return 0, err
			}
		}
	}
}

// prunes old commits and removes the least recently used repositories until the cache fits WasmGitCacheMaxSize,
// repositories in use are skipped
func (c *gitRepoCache) maintain(ctx context.Context) {
	entries, err := os.ReadDir(c.dir)
	if err != nil {
		return
	}

	type repoInfo struct {
		dir       string
		size      int64
		last_used time.Time
	}
	repos := []repoInfo{}
	var total int64
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		dir := filepath.Join(c.dir, e.Name())
		repo, ok := c.tryMaintain(dir)
		if ok {
			err = gitPruneCommits(ctx, dir, WasmGitCacheMaxCommits)
			if err != nil {
				logrus.WithError(err).WithField("dir", dir).Warn("git cache: prune failed")
			}
			c.doneMaintain(repo)
		}

		size, err := dirSize(dir)
		if err != nil {
			continue
		}
		last_used := time.Time{}
		if info, err := os.Stat(filepath.Join(dir, gitCacheUsedFile)); err == nil {
			last_used = info.ModTime()
		}
		repos = append(repos, repoInfo{dir: dir, size: size, last_used: last_used})
		total += size
	}

	sort.Slice(repos, func(i, j int) bool { return repos[i].last_used.Before(repos[j].last_used) })
	for _, r := range repos {
		if total <= WasmGitCacheMaxSize {
			break
		}
		repo, ok := c.tryMaintain(r.dir)
		if !ok {
			continue
		}
		logrus.WithFields(logrus.Fields{"dir": r.dir, "size": r.size}).Debug("git cache: evicting")
		err = os.RemoveAll(r.dir)
		c.doneMaintain(repo)
		if err != nil && !errors.Is(err, fs.ErrNotExist) {
			continue
		}
		total -= r.size
		metrics.WasmGitEvictions.Inc()
	}
	metrics.WasmGitCacheBytes.Set(float64(total))
}

// removes all but the newest (committer date) keep commits once there are more than 2*keep, the repository must not be in use
func gitPruneCommits(ctx context.Context, dir string, keep int) error {
	if _, err := os.Stat(filepath.Join(dir, "HEAD")); err != nil {
		return nil
	}
	out, err := runGit(ctx, dir, "for-each-ref", "--sort=-committerdate", "--format=%(refname)", gitCacheRefs)
	if err != nil {
		return err
	}
	refs := strings.Fields(out)
	if len(refs) <= 2*keep {
		return nil
	}
	stdin := &strings.Builder{}
	for _, ref := range refs[keep:] {
		fmt.Fprintf(stdin, "delete %s\n", ref)
	}
	_, err = runGitStdin(ctx, dir, stdin.String(), "update-ref", "--stdin")
	if err != nil {
		return err
	}
	_, err = runGit(ctx, dir, "gc", "--quiet", "--prune=now")
	return err
}
