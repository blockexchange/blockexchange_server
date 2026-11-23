package core

import (
	"blockexchange/types"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/sync/singleflight"
)

// disk cache for the builtin and dependency packs, keyed by the resolved commits of their entries

var (
	// total size of the cached packs, the least recently used ones are removed first
	WasmPackCacheMaxSize = int64(2 * 1024 * 1024 * 1024)
	// how long a resolved "latest", branch or tag commit is reused before asking the remote again
	WasmPackResolveTTL = 10 * time.Minute
	// max time for a pack build (keeps running if the requesting client goes away, to fill the cache)
	WasmPackBuildTimeout = 5 * time.Minute
)

// bump if the layout or content of the cached packs changes, invalidates all cached packs
const wasmPackFormat = "3"

// contents of a cached pack, needed for the config pack
type WasmPackManifest struct {
	// id of the game in the pack, empty if none
	GameID string `json:"gameid,omitempty"`
	// mod folder names
	Mods []string `json:"mods"`
}

type WasmPack struct {
	// cached pack file, don't modify
	Path string
	// http etag (quoted), changes with the content
	ETag     string
	ModTime  time.Time
	Manifest *WasmPackManifest
}

type wasmPackCache struct {
	dir string
	sf  singleflight.Group

	resolved_lock sync.Mutex
	// url + "\x00" + version -> resolved commit
	resolved map[string]resolvedCommit
}

type resolvedCommit struct {
	commit string
	at     time.Time
}

func newWasmPackCache(dir string) *wasmPackCache {
	if dir == "" {
		dir = filepath.Join(os.TempDir(), "blockexchange-wasm-packs")
	}
	return &wasmPackCache{dir: dir, resolved: map[string]resolvedCommit{}}
}

// returns the cached pack or builds it, concurrent requests for the same pack share a build
func (c *Core) GetWasmPack(ctx context.Context, kind WasmPackKind, deps []*types.SchemaModDependency) (*WasmPack, error) {
	return c.getWasmPack(ctx, kind, deps, nil)
}

// the builtin mods without the given mods (provided by the schematic), the full pack is shared if nothing overlaps
func (c *Core) GetWasmBuiltinPack(ctx context.Context, exclude_mods []string) (*WasmPack, error) {
	deps := WasmPackBuiltinDependencies()
	full, err := c.getWasmPack(ctx, WasmPackKindBuiltin, deps, nil)
	if err != nil {
		return nil, err
	}
	overlap := []string{}
	for _, name := range full.Manifest.Mods {
		if slices.Contains(exclude_mods, name) {
			overlap = append(overlap, name)
		}
	}
	if len(overlap) == 0 {
		return full, nil
	}
	return c.getWasmPack(ctx, WasmPackKindBuiltin, deps, overlap)
}

func (c *Core) getWasmPack(ctx context.Context, kind WasmPackKind, deps []*types.SchemaModDependency, exclude_mods []string) (*WasmPack, error) {
	if kind == WasmPackKindDependencies {
		game_count := 0
		for _, dep := range deps {
			if dep.ModType == types.SchemaModDependencyTypeGame {
				game_count++
			}
		}
		if game_count > 1 {
			return nil, invalidPack("only one game dependency is allowed, found %d", game_count)
		}
	}

	cache := c.wasmPacks
	resolved := make([]resolvedDependency, len(deps))
	for i, dep := range deps {
		commit, err := cache.resolveCommit(ctx, dep.GitURL, dep.Version)
		if err != nil {
			return nil, err
		}
		resolved[i] = resolvedDependency{dep: dep, commit: commit}
	}

	key := wasmPackKey(kind, resolved, exclude_mods)
	if p := cache.lookup(key); p != nil {
		return p, nil
	}

	ch := cache.sf.DoChan(key, func() (any, error) {
		if p := cache.lookup(key); p != nil {
			return p, nil
		}
		// detached from the request, the result is cached for the next one
		build_ctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), WasmPackBuildTimeout)
		defer cancel()
		return cache.build(build_ctx, key, resolved, exclude_mods)
	})
	select {
	case res := <-ch:
		if res.Err != nil {
			return nil, res.Err
		}
		return res.Val.(*WasmPack), nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func wasmPackKey(kind WasmPackKind, deps []resolvedDependency, exclude_mods []string) string {
	lines := make([]string, len(deps))
	for i, rd := range deps {
		lines[i] = strings.Join([]string{string(rd.dep.ModType), rd.dep.GitURL, rd.commit}, "\x00")
	}
	// the order of the entries doesn't change the content
	sort.Strings(lines)
	exclude := slices.Clone(exclude_mods)
	sort.Strings(exclude)
	h := sha256.New()
	fmt.Fprintf(h, "%s\n%s\n%s\nexclude:%s", wasmPackFormat, kind, strings.Join(lines, "\n"), strings.Join(exclude, ","))
	return string(kind) + "-" + hex.EncodeToString(h.Sum(nil))
}

func (cache *wasmPackCache) resolveCommit(ctx context.Context, url, version string) (string, error) {
	if gitCommitRegex.MatchString(version) {
		return strings.ToLower(version), nil
	}

	key := url + "\x00" + version
	cache.resolved_lock.Lock()
	rc, found := cache.resolved[key]
	cache.resolved_lock.Unlock()
	if found && time.Since(rc.at) < WasmPackResolveTTL {
		return rc.commit, nil
	}

	commit, err := gitResolveCommit(ctx, url, version)
	if err != nil {
		return "", err
	}
	cache.resolved_lock.Lock()
	cache.resolved[key] = resolvedCommit{commit: commit, at: time.Now()}
	cache.resolved_lock.Unlock()
	return commit, nil
}

func (cache *wasmPackCache) packPath(key string) string {
	return filepath.Join(cache.dir, key+".pack")
}

func (cache *wasmPackCache) manifestPath(key string) string {
	return filepath.Join(cache.dir, key+".json")
}

// returns the cached pack or nil
func (cache *wasmPackCache) lookup(key string) *WasmPack {
	data, err := os.ReadFile(cache.manifestPath(key))
	if err != nil {
		return nil
	}
	manifest := &WasmPackManifest{}
	if json.Unmarshal(data, manifest) != nil {
		return nil
	}
	pack_path := cache.packPath(key)
	info, err := os.Stat(pack_path)
	if err != nil {
		return nil
	}

	// last use for the eviction is the modtime of the manifest, the pack keeps the modtime of its content
	now := time.Now()
	os.Chtimes(cache.manifestPath(key), now, now)

	return &WasmPack{
		Path:     pack_path,
		ETag:     `"` + key + `"`,
		ModTime:  info.ModTime(),
		Manifest: manifest,
	}
}

func (cache *wasmPackCache) build(ctx context.Context, key string, deps []resolvedDependency, exclude_mods []string) (*WasmPack, error) {
	err := os.MkdirAll(cache.dir, 0755)
	if err != nil {
		return nil, err
	}

	// written to temp files first, the manifest is moved into place before the pack:
	// a pack file without manifest is never visible to lookup()
	f, err := os.CreateTemp(cache.dir, ".tmp-*.pack")
	if err != nil {
		return nil, err
	}
	defer os.Remove(f.Name())
	manifest, err := buildWasmPack(ctx, deps, exclude_mods, f)
	if err != nil {
		f.Close()
		return nil, err
	}
	err = f.Close()
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(manifest)
	if err != nil {
		return nil, err
	}
	err = os.WriteFile(cache.manifestPath(key), data, 0644)
	if err != nil {
		return nil, err
	}
	err = os.Rename(f.Name(), cache.packPath(key))
	if err != nil {
		return nil, err
	}

	cache.evict(key)

	p := cache.lookup(key)
	if p == nil {
		return nil, fmt.Errorf("pack '%s' missing after build", key)
	}
	return p, nil
}

// removes the least recently used packs until the cache fits WasmPackCacheMaxSize, keeps the given pack
func (cache *wasmPackCache) evict(keep string) {
	entries, err := os.ReadDir(cache.dir)
	if err != nil {
		return
	}

	type packFile struct {
		key       string
		size      int64
		last_used time.Time
	}
	packs := []packFile{}
	var total int64
	for _, e := range entries {
		info, err := e.Info()
		if err != nil {
			continue
		}
		name := e.Name()
		if strings.HasPrefix(name, ".tmp-") {
			// leftover of an interrupted build
			if time.Since(info.ModTime()) > 2*WasmPackBuildTimeout {
				os.Remove(filepath.Join(cache.dir, name))
			}
			continue
		}
		key, is_pack := strings.CutSuffix(name, ".pack")
		if !is_pack {
			continue
		}
		last_used := info.ModTime()
		if minfo, err := os.Stat(cache.manifestPath(key)); err == nil {
			last_used = minfo.ModTime()
		}
		packs = append(packs, packFile{key: key, size: info.Size(), last_used: last_used})
		total += info.Size()
	}

	sort.Slice(packs, func(i, j int) bool { return packs[i].last_used.Before(packs[j].last_used) })
	for _, p := range packs {
		if total <= WasmPackCacheMaxSize {
			break
		}
		if p.key == keep {
			continue
		}
		logrus.WithFields(logrus.Fields{"key": p.key, "size": p.size}).Debug("wasm pack cache: evicting")
		os.Remove(cache.packPath(p.key))
		os.Remove(cache.manifestPath(p.key))
		total -= p.size
	}
}
