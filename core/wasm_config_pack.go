package core

import (
	"archive/tar"
	"bytes"
	"fmt"
	"slices"
	"strings"

	"github.com/klauspost/compress/zstd"
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
secure.http_mods = blockexchange
`

const wasmPackWorldMT = `creative_mode = true
enable_damage = false
server_announce = false
backend = sqlite3
auth_backend = sqlite3
player_backend = sqlite3
mod_storage_backend = sqlite3
`

// world.mt and minetest.conf for the given dependency and builtin packs, small and generated per request
func (c *Core) WasmConfigPack(world_name string, deps, builtin *WasmPackManifest) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	zw, err := zstd.NewWriter(buf)
	if err != nil {
		return nil, err
	}
	tw := tar.NewWriter(zw)
	for _, d := range []string{"luanti", "luanti/worlds", WasmPackWorldDir} {
		err = tw.WriteHeader(&tar.Header{Typeflag: tar.TypeDir, Name: d + "/", Mode: 0755})
		if err != nil {
			return nil, err
		}
	}

	err = writeTarFile(tw, WasmPackMinetestConf, []byte(c.wasmMinetestConf()))
	if err != nil {
		return nil, err
	}
	err = writeTarFile(tw, WasmPackWorldDir+"/world.mt", []byte(wasmWorldMT(world_name, deps, builtin)))
	if err != nil {
		return nil, err
	}

	err = tw.Close()
	if err != nil {
		return nil, err
	}
	err = zw.Close()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// client settings, the blockexchange mod talks to this server
func (c *Core) wasmMinetestConf() string {
	sb := strings.Builder{}
	sb.WriteString(wasmPackMinetestConf)
	if c.cfg.BaseURL != "" {
		fmt.Fprintf(&sb, "blockexchange.url = %s\n", c.cfg.BaseURL)
	}
	return sb.String()
}

// world settings: the game and worldmods are always loaded, the builtin mods in luanti/mods
// are enabled unless a mod of the schematic has the same name
func wasmWorldMT(world_name string, deps, builtin *WasmPackManifest) string {
	sb := strings.Builder{}
	if deps.GameID != "" {
		// the game in the "game" folder of the world is used regardless of the id
		fmt.Fprintf(&sb, "gameid = %s\n", deps.GameID)
	}
	fmt.Fprintf(&sb, "world_name = %s\n", strings.NewReplacer("\n", " ", "\r", " ").Replace(world_name))
	sb.WriteString(wasmPackWorldMT)
	for _, name := range builtin.Mods {
		if !slices.Contains(deps.Mods, name) {
			fmt.Fprintf(&sb, "load_mod_%s = true\n", name)
		}
	}
	return sb.String()
}
