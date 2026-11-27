package core

import (
	"archive/tar"
	"blockexchange/types"
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"slices"
	"strings"

	"github.com/klauspost/compress/zstd"
)

// packs for the luanti wasm client: zstd compressed tar archives, extracted into the filesystem root of the client.
//
// The client installs three packs, split by how often they change so the heavy ones can be cached and shared:
//   - dependencies: the game and mods of a schematic in the world folder (shared by schematics with the same dependencies)
//   - builtin: the blockexchange mod and its dependencies in worldmods (see WasmPackBuiltinMods), without the mods
//     the schematic brings itself so the packs never write into the same folder (shared if there is no overlap)
//   - config: world.mt, bx_provision.json and minetest.conf, generated per request (see wasm_config_pack.go)

const (
	WasmPackWorldDir     = "luanti/worlds/world"
	WasmPackGameDir      = WasmPackWorldDir + "/game"
	WasmPackWorldModDir  = WasmPackWorldDir + "/worldmods"
	WasmPackMinetestConf = "luanti/minetest.conf"

	// upper limit of the uncompressed pack content
	WasmPackMaxSize = int64(1024 * 1024 * 1024)
)

type WasmPackKind string

const (
	// game -> world/game, mods -> world/worldmods
	WasmPackKindDependencies WasmPackKind = "dependencies"
	// mods -> world/worldmods
	WasmPackKindBuiltin WasmPackKind = "builtin"
)

// the blockexchange mod and its hard dependencies, a dependency of the schematic with the same mod name takes precedence
var WasmPackBuiltinMods = []string{
	"https://github.com/blockexchange/blockexchange",
	"https://github.com/buckaroobanzay/mtzip",
	"https://github.com/buckaroobanzay/placeholder",
	"https://github.com/mt-mods/promise",
	"https://github.com/Uberi/Minetest-WorldEdit",
	"https://github.com/minetest-mods/unified_inventory.git",
}

// transport protocols git may use for the dependency urls (no file://, ext:: etc)
var WasmPackGitProtocols = "http:https"

// validation and limit errors, caused by the dependency entries and not by the server
var ErrWasmPackInvalid = errors.New("invalid pack")

var gitCommitRegex = regexp.MustCompile(`^([0-9a-fA-F]{40}|[0-9a-fA-F]{64})$`)
var modNameRegex = regexp.MustCompile(`^[a-z0-9_]+$`)
var invalidModNameChars = regexp.MustCompile(`[^a-z0-9_]`)

// limits the number of concurrent pack builds (fetching is io and network heavy)
var wasmPackSemaphore = make(chan struct{}, 2)

func invalidPack(format string, args ...any) error {
	return fmt.Errorf("%w: %s", ErrWasmPackInvalid, fmt.Sprintf(format, args...))
}

// the builtin mods as dependency entries
func WasmPackBuiltinDependencies() []*types.SchemaModDependency {
	deps := make([]*types.SchemaModDependency, len(WasmPackBuiltinMods))
	for i, url := range WasmPackBuiltinMods {
		deps[i] = &types.SchemaModDependency{
			ModType: types.SchemaModDependencyTypeMod,
			GitURL:  url,
			Version: types.SchemaModDependencyDefaultVersion,
		}
	}
	return deps
}

// a dependency with its version resolved to a commit
type resolvedDependency struct {
	dep    *types.SchemaModDependency
	commit string
}

type wasmPackEntry struct {
	archive_dir string
	repo        *gitRepo
	commit      string
}

// fetches the resolved dependencies into the git cache and writes the pack and its manifest, mods with a name in exclude_mods are left out
func buildWasmPack(ctx context.Context, git *gitRepoCache, deps []resolvedDependency, exclude_mods []string, w io.Writer) (*WasmPackManifest, error) {
	select {
	case wasmPackSemaphore <- struct{}{}:
		defer func() { <-wasmPackSemaphore }()
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	manifest := &WasmPackManifest{Mods: []string{}}
	mod_urls := map[string]string{}
	entries := []wasmPackEntry{}

	// everything is fetched before the first byte is written
	for _, rd := range deps {
		repo := git.acquire(rd.dep.GitURL)
		defer git.release(repo)
		err := git.ensureCommit(ctx, repo, rd.dep.GitURL, rd.commit)
		if err != nil {
			return nil, err
		}

		if rd.dep.ModType == types.SchemaModDependencyTypeGame {
			manifest.GameID = repoName(rd.dep.GitURL)
			entries = append(entries, wasmPackEntry{archive_dir: WasmPackGameDir, repo: repo, commit: rd.commit})
			continue
		}

		name := modFolderName(ctx, repo, rd.commit, rd.dep.GitURL)
		if slices.Contains(exclude_mods, name) {
			continue
		}
		if other_url := mod_urls[name]; other_url != "" {
			return nil, invalidPack("'%s' and '%s' both resolve to the mod folder '%s'", other_url, rd.dep.GitURL, name)
		}
		mod_urls[name] = rd.dep.GitURL
		manifest.Mods = append(manifest.Mods, name)

		entries = append(entries, wasmPackEntry{archive_dir: WasmPackWorldModDir + "/" + name, repo: repo, commit: rd.commit})
	}

	skeleton := []string{"luanti", "luanti/worlds", WasmPackWorldDir, WasmPackWorldModDir}

	zw, err := zstd.NewWriter(w)
	if err != nil {
		return nil, err
	}
	tw := tar.NewWriter(zw)
	for _, d := range skeleton {
		err = tw.WriteHeader(&tar.Header{Typeflag: tar.TypeDir, Name: d + "/", Mode: 0755})
		if err != nil {
			return nil, err
		}
	}
	budget := WasmPackMaxSize
	for _, e := range entries {
		size, err := gitWriteTar(ctx, tw, e.repo, e.commit, e.archive_dir, budget)
		if err != nil {
			return nil, err
		}
		budget -= size
	}
	err = tw.Close()
	if err != nil {
		return nil, err
	}
	return manifest, zw.Close()
}

func writeTarFile(tw *tar.Writer, name string, data []byte) error {
	err := tw.WriteHeader(&tar.Header{Typeflag: tar.TypeReg, Name: name, Mode: 0644, Size: int64(len(data))})
	if err != nil {
		return err
	}
	_, err = tw.Write(data)
	return err
}

func dirSize(dir string) (int64, error) {
	var size int64
	err := filepath.WalkDir(dir, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() && d.Name() == ".git" {
			return filepath.SkipDir
		}
		if d.Type().IsRegular() {
			info, err := d.Info()
			if err != nil {
				return err
			}
			size += info.Size()
		}
		return nil
	})
	return size, err
}

// resolves "latest" (default branch), a branch or a tag to a commit with "git ls-remote", full commit hashes are returned as-is
func gitResolveCommit(ctx context.Context, url, version string) (string, error) {
	if version == "" {
		version = types.SchemaModDependencyDefaultVersion
	}
	if gitCommitRegex.MatchString(version) {
		return strings.ToLower(version), nil
	}

	pattern := "HEAD"
	if version != types.SchemaModDependencyDefaultVersion {
		if !types.GitVersionRegex.MatchString(version) {
			return "", invalidPack("invalid version '%s'", version)
		}
		pattern = version
	}

	out, err := runGit(ctx, "", "ls-remote", "--end-of-options", url, pattern)
	if err != nil {
		return "", err
	}
	refs := map[string]string{}
	for _, line := range strings.Split(out, "\n") {
		commit, ref, found := strings.Cut(strings.TrimSpace(line), "\t")
		if found {
			refs[ref] = commit
		}
	}

	candidates := []string{"HEAD"}
	if version != types.SchemaModDependencyDefaultVersion {
		// annotated tags point to the tag object, "^{}" is the commit
		candidates = []string{"refs/tags/" + version + "^{}", "refs/tags/" + version, "refs/heads/" + version, version}
	}
	for _, c := range candidates {
		if commit := refs[c]; gitCommitRegex.MatchString(commit) {
			return strings.ToLower(commit), nil
		}
	}
	return "", invalidPack("version '%s' not found in '%s' (use 'latest', a branch, a tag or a full commit hash)", version, url)
}

func gitEnv() []string {
	return append(os.Environ(),
		"GIT_TERMINAL_PROMPT=0",
		"GIT_ASKPASS=/bin/true",
		"GIT_ALLOW_PROTOCOL="+WasmPackGitProtocols,
		"GIT_CONFIG_NOSYSTEM=1",
		"GIT_CONFIG_GLOBAL=/dev/null",
		"GIT_LFS_SKIP_SMUDGE=1",
	)
}

func runGit(ctx context.Context, dir string, args ...string) (string, error) {
	return runGitStdin(ctx, dir, "", args...)
}

func runGitStdin(ctx context.Context, dir, stdin string, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, "git", args...)
	cmd.Dir = dir
	cmd.Env = gitEnv()
	if stdin != "" {
		cmd.Stdin = strings.NewReader(stdin)
	}
	out, err := cmd.CombinedOutput()
	if ctx.Err() != nil {
		return "", ctx.Err()
	}
	if err != nil {
		return "", fmt.Errorf("git %s failed: %s", args[0], strings.TrimSpace(string(out)))
	}
	return string(out), nil
}

// folder name of a mod: "name" in mod.conf or modpack.conf, the repository name otherwise
func modFolderName(ctx context.Context, repo *gitRepo, commit, url string) string {
	for _, conf := range []string{"mod.conf", "modpack.conf"} {
		name := gitReadConfValue(ctx, repo, commit, conf, "name")
		if modNameRegex.MatchString(name) {
			return name
		}
	}
	return repoName(url)
}

// lowercase repository name with only a-z, 0-9 and _ ("https://x/Minetest-WorldEdit.git" -> "minetest_worldedit")
func repoName(url string) string {
	name := strings.TrimSuffix(path.Base(strings.TrimRight(url, "/")), ".git")
	name = invalidModNameChars.ReplaceAllString(strings.ToLower(name), "_")
	if name == "" {
		name = "mod"
	}
	return name
}

func readConfValueFrom(r io.Reader, key string) string {
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		k, v, found := strings.Cut(sc.Text(), "=")
		if found && strings.TrimSpace(k) == key {
			return strings.TrimSpace(v)
		}
	}
	return ""
}
