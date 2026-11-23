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
	"strings"

	"github.com/klauspost/compress/zstd"
)

const (
	// archive paths, the luanti wasm client extracts the pack into its filesystem root
	WasmPackWorldDir     = "luanti/worlds/world"
	WasmPackGameDir      = WasmPackWorldDir + "/game"
	WasmPackWorldModDir  = WasmPackWorldDir + "/worldmods"
	WasmPackMinetestConf = "luanti/minetest.conf"

	// upper limit of the uncompressed pack content
	WasmPackMaxSize = int64(1024 * 1024 * 1024)
)

// client settings: creative building with fly/fast/noclip
const wasmPackMinetestConf = `creative_mode = true
enable_damage = false
default_privs = interact, shout, fly, fast, noclip, give
fast_move = true
free_move = true
noclip = true
server_announce = false
mg_name = singlenode
`

const wasmPackWorldMT = `creative_mode = true
enable_damage = false
server_announce = false
backend = sqlite3
auth_backend = sqlite3
player_backend = sqlite3
mod_storage_backend = sqlite3
`

// mods that are always added to the pack: the blockexchange mod and its hard dependencies,
// skipped if a dependency of the schematic already provides a mod with the same name
var WasmPackBuiltinMods = []string{
	"https://github.com/blockexchange/blockexchange",
	"https://github.com/buckaroobanzay/mtzip",
	"https://github.com/buckaroobanzay/placeholder",
	"https://github.com/mt-mods/promise",
}

// transport protocols git may use for the dependency urls (no file://, ext:: etc)
var WasmPackGitProtocols = "http:https"

// validation and limit errors, caused by the dependency entries and not by the server
var ErrWasmPackInvalid = errors.New("invalid pack")

// branch, tag or full commit hash, must not look like a command line option
var GitVersionRegex = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._/-]*$`)
var modNameRegex = regexp.MustCompile(`^[a-z0-9_]+$`)
var invalidModNameChars = regexp.MustCompile(`[^a-z0-9_]`)

// limits the number of concurrent pack builds (cloning is io and network heavy)
var wasmPackSemaphore = make(chan struct{}, 2)

type WasmPack struct {
	dir        string
	world_name string
	// blockexchange.url setting, the url of this server
	server_url string
	// id of the embedded game, empty if there is no game dependency
	gameid string
	// archive dir -> local dir
	entries []wasmPackEntry
}

type wasmPackEntry struct {
	archive_dir string
	local_dir   string
}

func invalidPack(format string, args ...any) error {
	return fmt.Errorf("%w: %s", ErrWasmPackInvalid, fmt.Sprintf(format, args...))
}

// clones all dependencies into a temporary directory, Close() removes it
func (c *Core) PrepareWasmPack(ctx context.Context, world_name string, deps []*types.SchemaModDependency) (*WasmPack, error) {
	select {
	case wasmPackSemaphore <- struct{}{}:
		defer func() { <-wasmPackSemaphore }()
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	game_count := 0
	for _, dep := range deps {
		if dep.ModType == types.SchemaModDependencyTypeGame {
			game_count++
		}
	}
	if game_count > 1 {
		return nil, invalidPack("only one game dependency is allowed, found %d", game_count)
	}

	dir, err := os.MkdirTemp("", "wasm-pack-")
	if err != nil {
		return nil, fmt.Errorf("create temp dir: %w", err)
	}
	pack := &WasmPack{dir: dir, world_name: world_name, server_url: c.cfg.BaseURL}

	mod_names := map[string]string{}
	var total_size int64
	add := func(i int, dep *types.SchemaModDependency, builtin bool) error {
		local_dir := filepath.Join(dir, fmt.Sprintf("%d", i))
		err := gitShallowClone(ctx, dep.GitURL, dep.Version, local_dir)
		if err != nil {
			return err
		}

		archive_dir := WasmPackGameDir
		if dep.ModType == types.SchemaModDependencyTypeGame {
			pack.gameid = repoName(dep.GitURL)
		}
		if dep.ModType == types.SchemaModDependencyTypeMod {
			name := modFolderName(local_dir, dep.GitURL)
			if other_url := mod_names[name]; other_url != "" {
				if builtin {
					// provided by the schematic, possibly with another version
					return os.RemoveAll(local_dir)
				}
				return invalidPack("'%s' and '%s' both resolve to the mod folder '%s'", other_url, dep.GitURL, name)
			}
			mod_names[name] = dep.GitURL
			archive_dir = WasmPackWorldModDir + "/" + name
		}

		size, err := dirSize(local_dir)
		if err != nil {
			return err
		}
		total_size += size
		if total_size > WasmPackMaxSize {
			return invalidPack("pack exceeds the size limit of %d bytes", WasmPackMaxSize)
		}

		pack.entries = append(pack.entries, wasmPackEntry{archive_dir: archive_dir, local_dir: local_dir})
		return nil
	}

	for i, dep := range deps {
		err = add(i, dep, false)
		if err != nil {
			pack.Close()
			return nil, err
		}
	}
	for i, url := range WasmPackBuiltinMods {
		dep := &types.SchemaModDependency{
			ModType: types.SchemaModDependencyTypeMod,
			GitURL:  url,
			Version: types.SchemaModDependencyDefaultVersion,
		}
		err = add(len(deps)+i, dep, true)
		if err != nil {
			pack.Close()
			return nil, err
		}
	}

	return pack, nil
}

func (p *WasmPack) Close() error {
	return os.RemoveAll(p.dir)
}

// writes the zstd compressed tar archive
func (p *WasmPack) Write(w io.Writer) error {
	zw, err := zstd.NewWriter(w)
	if err != nil {
		return err
	}
	tw := tar.NewWriter(zw)

	for _, d := range []string{"luanti", "luanti/worlds", WasmPackWorldDir, WasmPackWorldModDir} {
		err = tw.WriteHeader(&tar.Header{Typeflag: tar.TypeDir, Name: d + "/", Mode: 0755})
		if err != nil {
			return err
		}
	}

	err = writeTarFile(tw, WasmPackMinetestConf, []byte(p.minetestConf()))
	if err != nil {
		return err
	}
	err = writeTarFile(tw, WasmPackWorldDir+"/world.mt", []byte(p.worldMT()))
	if err != nil {
		return err
	}

	for _, e := range p.entries {
		err = writeTarDir(tw, e.local_dir, e.archive_dir)
		if err != nil {
			return err
		}
	}

	err = tw.Close()
	if err != nil {
		return err
	}
	return zw.Close()
}

// client settings, the blockexchange mod talks to this server
func (p *WasmPack) minetestConf() string {
	sb := strings.Builder{}
	sb.WriteString(wasmPackMinetestConf)
	sb.WriteString("secure.http_mods = blockexchange\n")
	if p.server_url != "" {
		fmt.Fprintf(&sb, "blockexchange.url = %s\n", p.server_url)
	}
	return sb.String()
}

// world settings, the mods in "worldmods" are always enabled and need no load_mod_* entry
func (p *WasmPack) worldMT() string {
	sb := strings.Builder{}
	if p.gameid != "" {
		// the game in the "game" folder of the world is used regardless of the id
		fmt.Fprintf(&sb, "gameid = %s\n", p.gameid)
	}
	fmt.Fprintf(&sb, "world_name = %s\n", strings.NewReplacer("\n", " ", "\r", " ").Replace(p.world_name))
	sb.WriteString(wasmPackWorldMT)
	return sb.String()
}

func writeTarFile(tw *tar.Writer, name string, data []byte) error {
	err := tw.WriteHeader(&tar.Header{Typeflag: tar.TypeReg, Name: name, Mode: 0644, Size: int64(len(data))})
	if err != nil {
		return err
	}
	_, err = tw.Write(data)
	return err
}

// adds the directory recursively, skips the .git folder, symlinks and special files
func writeTarDir(tw *tar.Writer, local_dir, archive_dir string) error {
	return filepath.WalkDir(local_dir, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(local_dir, p)
		if err != nil {
			return err
		}
		if d.IsDir() && d.Name() == ".git" {
			return filepath.SkipDir
		}
		if !d.IsDir() && !d.Type().IsRegular() {
			// symlinks could point outside of the repository
			return nil
		}

		info, err := d.Info()
		if err != nil {
			return err
		}
		hdr, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}
		hdr.Name = path.Join(archive_dir, filepath.ToSlash(rel))
		hdr.Uid, hdr.Gid, hdr.Uname, hdr.Gname = 0, 0, "", ""
		if d.IsDir() {
			hdr.Name += "/"
			return tw.WriteHeader(hdr)
		}

		err = tw.WriteHeader(hdr)
		if err != nil {
			return err
		}
		f, err := os.Open(p)
		if err != nil {
			return err
		}
		defer f.Close()
		_, err = io.Copy(tw, f)
		return err
	})
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

// shallow clone of the default branch ("latest") or the given branch, tag or full commit hash
func gitShallowClone(ctx context.Context, url, version, dir string) error {
	if version == "" || version == types.SchemaModDependencyDefaultVersion {
		return runGit(ctx, "", "clone", "--quiet", "--depth", "1", "--single-branch", "--no-tags", "--end-of-options", url, dir)
	}

	if !GitVersionRegex.MatchString(version) {
		return invalidPack("invalid version '%s'", version)
	}
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}
	err = runGit(ctx, dir, "init", "--quiet")
	if err != nil {
		return err
	}
	err = runGit(ctx, dir, "fetch", "--quiet", "--depth", "1", "--no-tags", "--end-of-options", url, version)
	if err != nil {
		return err
	}
	return runGit(ctx, dir, "checkout", "--quiet", "--detach", "FETCH_HEAD")
}

func runGit(ctx context.Context, dir string, args ...string) error {
	cmd := exec.CommandContext(ctx, "git", args...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(),
		"GIT_TERMINAL_PROMPT=0",
		"GIT_ASKPASS=/bin/true",
		"GIT_ALLOW_PROTOCOL="+WasmPackGitProtocols,
		"GIT_CONFIG_NOSYSTEM=1",
		"GIT_CONFIG_GLOBAL=/dev/null",
		"GIT_LFS_SKIP_SMUDGE=1",
	)
	out, err := cmd.CombinedOutput()
	if ctx.Err() != nil {
		return ctx.Err()
	}
	if err != nil {
		return fmt.Errorf("git %s failed: %s", args[0], strings.TrimSpace(string(out)))
	}
	return nil
}

// folder name of a mod: "name" in mod.conf or modpack.conf, the repository name otherwise
func modFolderName(dir, url string) string {
	for _, conf := range []string{"mod.conf", "modpack.conf"} {
		name := readConfValue(filepath.Join(dir, conf), "name")
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

func readConfValue(file, key string) string {
	info, err := os.Lstat(file)
	if err != nil || !info.Mode().IsRegular() {
		return ""
	}
	f, err := os.Open(file)
	if err != nil {
		return ""
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		k, v, found := strings.Cut(sc.Text(), "=")
		if found && strings.TrimSpace(k) == key {
			return strings.TrimSpace(v)
		}
	}
	return ""
}
