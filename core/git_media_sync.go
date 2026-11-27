package core

import (
	"blockexchange/types"
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"cirello.io/pglock"
	"github.com/sirupsen/logrus"
)

// imports the textures and models of git repositories into the media files:
// one bare repository per url in the cache dir, only the latest commit of the default branch is fetched (shallow)
// and the files are read from the object store without a checkout
//
//	<dir>/<sha256(url)>/        bare repository
//	<dir>/<sha256(url)>/bx-url  the url, for debugging

// upper bound for a whole sync run (fetch and import)
const gitMediaSyncTimeout = 15 * time.Minute

const (
	gitMediaRef = "refs/bx/head"
	// larger files are skipped
	gitMediaMaxFileSize = int64(16 * 1024 * 1024)
	// media files are stored once this many bytes are collected
	gitMediaBatchSize = int64(32 * 1024 * 1024)
)

// folders the media files are read from, at any depth so modpacks and games are covered too
var gitMediaFolders = map[string]bool{"textures": true, "models": true}

var ErrGitMediaSyncRunning = errors.New("a sync for this source is already running")

func (c *Core) gitMediaCacheDir() string {
	if c.cfg.GitMediaCacheDir != "" {
		return c.cfg.GitMediaCacheDir
	}
	return filepath.Join(os.TempDir(), "blockexchange-git-media")
}

type gitMediaEntry struct {
	name   string // file name ("default_stone.png")
	object string
	size   int64
}

// true if the file is a media file in one of the media folders ("textures/x.png", "mymod/models/sub/y.obj")
func isGitMediaPath(p string) bool {
	if !luantiMediaExtensions[strings.ToLower(path.Ext(p))] {
		return false
	}
	dirs := strings.Split(path.Dir(p), "/")
	for _, d := range dirs {
		if gitMediaFolders[d] {
			return true
		}
	}
	return false
}

// parses the output of "git ls-tree -r -l -z", returns the media files and the number of skipped (too large) files,
// the first file wins if the same name is in multiple folders
func parseGitMediaTree(out string) ([]gitMediaEntry, int) {
	entries := []gitMediaEntry{}
	seen := map[string]bool{}
	skipped := 0
	for _, line := range strings.Split(out, "\x00") {
		// "<mode> <type> <object> <size>\t<path>"
		meta, p, found := strings.Cut(line, "\t")
		fields := strings.Fields(meta)
		if !found || len(fields) != 4 || fields[1] != "blob" || !strings.HasPrefix(fields[0], "100") {
			// symlinks, submodules
			continue
		}
		if !isGitMediaPath(p) {
			continue
		}
		size, err := strconv.ParseInt(fields[3], 10, 64)
		if err != nil {
			continue
		}
		name := path.Base(p)
		if seen[name] {
			continue
		}
		seen[name] = true
		if size > gitMediaMaxFileSize || len(name) > 256 {
			skipped++
			continue
		}
		entries = append(entries, gitMediaEntry{name: name, object: fields[2], size: size})
	}
	return entries, skipped
}

// creates or updates the local copy of the repository, returns its directory
func gitMediaFetch(ctx context.Context, cache_dir, url string) (string, error) {
	dir := filepath.Join(cache_dir, gitRepoKey(url))
	if _, err := os.Stat(filepath.Join(dir, "HEAD")); err != nil {
		err = gitInitCacheRepo(ctx, dir, url)
		if err != nil {
			os.RemoveAll(dir)
			return "", err
		}
	}

	_, err := runGit(ctx, dir, "fetch", "--quiet", "--depth", "1", "--no-tags", "--no-write-fetch-head",
		"--end-of-options", url, "+HEAD:"+gitMediaRef)
	if err != nil {
		return "", err
	}
	// drops the objects of the previous commits
	_, err = runGit(ctx, dir, "gc", "--quiet", "--prune=now")
	if err != nil {
		return "", err
	}
	return dir, nil
}

// reads the blobs with a single "git cat-file --batch" process, fn is called in the order of the entries
func gitReadBlobs(ctx context.Context, dir string, entries []gitMediaEntry, fn func(e gitMediaEntry, data []byte) error) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	cmd := exec.CommandContext(ctx, "git", "cat-file", "--batch")
	cmd.Dir = dir
	cmd.Env = gitEnv()
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	err = cmd.Start()
	if err != nil {
		return err
	}
	go func() {
		w := bufio.NewWriter(stdin)
		for _, e := range entries {
			fmt.Fprintln(w, e.object)
		}
		w.Flush()
		stdin.Close()
	}()

	read_err := func() error {
		r := bufio.NewReader(stdout)
		for _, e := range entries {
			// "<object> <type> <size>\n<data>\n"
			header, err := r.ReadString('\n')
			if err != nil {
				return err
			}
			fields := strings.Fields(header)
			if len(fields) != 3 || fields[1] != "blob" {
				return fmt.Errorf("unexpected cat-file output for '%s': %s", e.name, strings.TrimSpace(header))
			}
			size, err := strconv.ParseInt(fields[2], 10, 64)
			if err != nil || size != e.size {
				return fmt.Errorf("unexpected size for '%s': %s", e.name, fields[2])
			}
			data := make([]byte, size+1)
			_, err = io.ReadFull(r, data)
			if err != nil {
				return err
			}
			err = fn(e, data[:size])
			if err != nil {
				return err
			}
		}
		return nil
	}()
	if read_err != nil {
		cancel()
		cmd.Wait()
		return read_err
	}
	return cmd.Wait()
}

// fills the empty license fields from the contentdb, the contentdb package is looked up by the git url if not set
func (c *Core) updateGitMediaSourceLicenses(ctx context.Context, src *types.GitMediaSource) error {
	if src.MediaLicense != "" && src.License != "" {
		return nil
	}
	if src.CDBPackage == "" {
		pkg, err := lookupCDBPackageByRepo(src.GitURL)
		if err != nil {
			return fmt.Errorf("contentdb package lookup: %w", err)
		}
		if pkg == "" {
			return errors.New("no contentdb package found for the git url")
		}
		src.CDBPackage = pkg
	}
	l, err := fetchCDBPackageLicenses(ctx, src.CDBPackage)
	if err != nil {
		return fmt.Errorf("contentdb package '%s': %w", src.CDBPackage, err)
	}
	if src.MediaLicense == "" {
		src.MediaLicense = l.MediaLicense
	}
	if src.License == "" {
		src.License = l.License
	}
	return c.repos.GitMediaSourceRepo.UpdateLicenses(src.ID, src.CDBPackage, src.MediaLicense, src.License)
}

func (c *Core) syncGitMediaSource(src *types.GitMediaSource) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), gitMediaSyncTimeout)
	defer cancel()

	dir, err := gitMediaFetch(ctx, c.gitMediaCacheDir(), src.GitURL)
	if err != nil {
		return "", fmt.Errorf("fetch: %w", err)
	}
	out, err := runGit(ctx, dir, "ls-tree", "-r", "-l", "-z", "--end-of-options", gitMediaRef)
	if err != nil {
		return "", fmt.Errorf("list files: %w", err)
	}
	entries, skipped := parseGitMediaTree(out)

	// store in batches to bound the memory usage
	batch := []*types.Mediafile{}
	var batch_size, total_size int64
	flush := func() error {
		err := c.repos.MediaRepo.UpsertMediafiles(batch)
		batch, batch_size = []*types.Mediafile{}, 0
		return err
	}
	err = gitReadBlobs(ctx, dir, entries, func(e gitMediaEntry, data []byte) error {
		batch = append(batch, &types.Mediafile{Name: e.name, GitMediaSourceID: &src.ID, Data: data})
		batch_size += int64(len(data))
		total_size += int64(len(data))
		if batch_size >= gitMediaBatchSize {
			return flush()
		}
		return nil
	})
	if err == nil {
		err = flush()
	}
	if err != nil {
		return "", fmt.Errorf("store media files: %w", err)
	}

	// files that were removed from the repository
	existing, err := c.repos.MediaRepo.GetMediafileNamesBySource(src.ID)
	if err != nil {
		return "", fmt.Errorf("list media files: %w", err)
	}
	current := make(map[string]bool, len(entries))
	for _, e := range entries {
		current[e.name] = true
	}
	removed := []string{}
	for _, name := range existing {
		if !current[name] {
			removed = append(removed, name)
		}
	}
	err = c.repos.MediaRepo.RemoveMediafilesOfSource(src.ID, removed)
	if err != nil {
		return "", fmt.Errorf("remove media files: %w", err)
	}

	msg := fmt.Sprintf("imported %d media files (%.1f MB)", len(entries), float64(total_size)/1024/1024)
	if len(removed) > 0 {
		msg += fmt.Sprintf(", removed %d", len(removed))
	}
	if skipped > 0 {
		msg += fmt.Sprintf(", skipped %d files larger than %d MB", skipped, gitMediaMaxFileSize/1024/1024)
	}

	// the media files are imported even if the licenses can't be determined
	err = c.updateGitMediaSourceLicenses(ctx, src)
	if err != nil {
		msg += fmt.Sprintf(", license lookup failed: %s", err.Error())
	}
	return msg, nil
}

// imports the media files of the source and writes the result back to it
func (c *Core) SyncGitMediaSource(src *types.GitMediaSource) error {
	lock, err := c.repos.PGLock.Acquire("git-media-sync-"+src.ID, pglock.FailIfLocked())
	if errors.Is(err, pglock.ErrNotAcquired) {
		return ErrGitMediaSyncRunning
	}
	if err != nil {
		return fmt.Errorf("acquire lock: %w", err)
	}
	defer lock.Close()

	log := logrus.WithFields(logrus.Fields{"source": src.ID, "git_url": src.GitURL})
	log.Info("syncing git media source")

	msg, syncErr := c.syncGitMediaSource(src)
	if syncErr != nil {
		msg = syncErr.Error()
		log.WithError(syncErr).Warn("git media sync failed")
	} else {
		log.Info(msg)
	}

	src.LastUpdatedTimestamp = time.Now().UnixMilli()
	src.LastUpdatedSuccess = syncErr == nil
	src.LastUpdatedMessage = msg
	return c.repos.GitMediaSourceRepo.UpdateStatus(src.ID, src.LastUpdatedTimestamp, src.LastUpdatedSuccess, src.LastUpdatedMessage)
}

// syncs every source whose last run is older than the given interval
func (c *Core) SyncDueGitMediaSources(interval time.Duration) {
	list, err := c.repos.GitMediaSourceRepo.GetAll()
	if err != nil {
		logrus.WithError(err).Error("list git media sources")
		return
	}
	for _, src := range list {
		if time.Since(time.UnixMilli(src.LastUpdatedTimestamp)) < interval {
			continue
		}
		err = c.SyncGitMediaSource(src)
		if err != nil && !errors.Is(err, ErrGitMediaSyncRunning) {
			logrus.WithError(err).WithField("source", src.ID).Error("git media sync")
		}
	}
}
