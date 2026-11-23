package core

import (
	"archive/tar"
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/klauspost/compress/zstd"
)

// static parts of the config pack, completed per request (see wasmMinetestConf and wasmWorldMT)

// client settings: creative building with fly/fast/noclip
//
//go:embed wasm_config/minetest.conf
var wasmPackMinetestConf string

// world settings
//
//go:embed wasm_config/world.mt
var wasmPackWorldMT string

// map settings of the world
//
//go:embed wasm_config/map_meta.txt
var wasmPackMapMeta string

// position of the provisioned schematic in the world
type WasmProvisionPosition struct {
	X int `json:"x"`
	Y int `json:"y"`
	Z int `json:"z"`
}

// bx_provision.json in the world folder: the schematic to place on startup
type WasmProvision struct {
	SchemaUID string                `json:"schema_uid"`
	Position  WasmProvisionPosition `json:"position"`
	// raw jwt of the user that loads the config pack
	Token string `json:"token"`
}

// world.mt, map_meta.txt, bx_provision.json and minetest.conf for the given schematic, user token and dependency pack, small and generated per request
func (c *Core) WasmConfigPack(schema_uid, token string, deps *WasmPackManifest) ([]byte, error) {
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
	err = writeTarFile(tw, WasmPackWorldDir+"/world.mt", []byte(wasmWorldMT(deps)))
	if err != nil {
		return nil, err
	}
	err = writeTarFile(tw, WasmPackWorldDir+"/map_meta.txt", []byte(wasmPackMapMeta))
	if err != nil {
		return nil, err
	}
	provision, err := json.MarshalIndent(&WasmProvision{SchemaUID: schema_uid, Token: token}, "", "\t")
	if err != nil {
		return nil, err
	}
	err = writeTarFile(tw, WasmPackWorldDir+"/bx_provision.json", provision)
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

// world settings: the game and all mods (worldmods) are always loaded, no load_mod_* entries needed
func wasmWorldMT(deps *WasmPackManifest) string {
	sb := strings.Builder{}
	if deps.GameID != "" {
		// the game in the "game" folder of the world is used regardless of the id
		fmt.Fprintf(&sb, "gameid = %s\n", deps.GameID)
	}
	sb.WriteString(wasmPackWorldMT)
	return sb.String()
}
