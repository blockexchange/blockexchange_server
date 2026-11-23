package api

import (
	"blockexchange/core"
	"blockexchange/types"
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

// packs for the luanti wasm client, see core/wasm_pack.go

// max time to wait for a pack (resolving and cloning)
const wasmPackRequestTimeout = 5 * time.Minute

func sendWasmPackError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, core.ErrWasmPackInvalid):
		SendError(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, context.DeadlineExceeded):
		SendError(w, http.StatusGatewayTimeout, "preparing the pack took too long")
	default:
		// git failures (unreachable repository, unknown commit)
		SendError(w, http.StatusBadGateway, err.Error())
	}
}

// serves a cached pack, the browser revalidates with the etag (304 if unchanged)
func serveWasmPack(w http.ResponseWriter, r *http.Request, pack *core.WasmPack, filename string) {
	f, err := os.Open(pack.Path)
	if err != nil {
		// evicted in the meantime
		SendError(w, http.StatusServiceUnavailable, "pack not available, please retry")
		return
	}
	defer f.Close()

	setWasmPackHeaders(w, pack.ETag, filename)
	http.ServeContent(w, r, filename, pack.ModTime, f)
}

func setWasmPackHeaders(w http.ResponseWriter, etag, filename string) {
	h := w.Header()
	h.Set("Content-Type", "application/zstd")
	h.Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	h.Set("ETag", etag)
	// per user (authenticated), always revalidate
	h.Set("Cache-Control", "private, no-cache")
}

// loads the dependencies of an owned schematic, nil if a response was sent
func (api *Api) loadWasmPackDependencies(w http.ResponseWriter, ctx *SecureContext, schema_uid string) (*types.Schema, []*types.SchemaModDependency) {
	schema := api.loadOwnedSchema(w, ctx, schema_uid, true)
	if schema == nil {
		return nil, nil
	}
	deps, err := api.SchemaModDependencyRepo.GetBySchemaUID(schema_uid)
	if SendIfError(w, err) {
		return nil, nil
	}
	if len(deps) == 0 {
		SendError(w, http.StatusNotFound, "the schematic has no mod dependencies")
		return nil, nil
	}
	return schema, deps
}

// the game and mods of the schematic (owner or admin), cached and shared between schematics with the same dependencies
func (api *Api) GetWasmPack(w http.ResponseWriter, r *http.Request, ctx *SecureContext) {
	schema, deps := api.loadWasmPackDependencies(w, ctx, mux.Vars(r)["schema_uid"])
	if schema == nil {
		return
	}

	req_ctx, cancel := context.WithTimeout(r.Context(), wasmPackRequestTimeout)
	defer cancel()
	pack, err := api.core.GetWasmPack(req_ctx, core.WasmPackKindDependencies, deps)
	if err != nil {
		sendWasmPackError(w, err)
		return
	}
	serveWasmPack(w, r, pack, schema.Name+".pack")
}

// the blockexchange mod and its dependencies (owner or admin), without the mods the schematic brings itself
func (api *Api) GetWasmBuiltinPack(w http.ResponseWriter, r *http.Request, ctx *SecureContext) {
	schema, deps := api.loadWasmPackDependencies(w, ctx, mux.Vars(r)["schema_uid"])
	if schema == nil {
		return
	}

	// usually a cache hit, the client installs the dependency pack first
	req_ctx, cancel := context.WithTimeout(r.Context(), wasmPackRequestTimeout)
	defer cancel()
	deps_pack, err := api.core.GetWasmPack(req_ctx, core.WasmPackKindDependencies, deps)
	if err != nil {
		sendWasmPackError(w, err)
		return
	}
	pack, err := api.core.GetWasmBuiltinPack(req_ctx, deps_pack.Manifest.Mods)
	if err != nil {
		sendWasmPackError(w, err)
		return
	}
	serveWasmPack(w, r, pack, schema.Name+".builtin.pack")
}

// world.mt, map_meta.txt and minetest.conf of the schematic (owner or admin), matching the dependency pack
func (api *Api) GetWasmConfigPack(w http.ResponseWriter, r *http.Request, ctx *SecureContext) {
	schema, deps := api.loadWasmPackDependencies(w, ctx, mux.Vars(r)["schema_uid"])
	if schema == nil {
		return
	}

	// usually a cache hit, the client installs the dependency pack first
	req_ctx, cancel := context.WithTimeout(r.Context(), wasmPackRequestTimeout)
	defer cancel()
	deps_pack, err := api.core.GetWasmPack(req_ctx, core.WasmPackKindDependencies, deps)
	if err != nil {
		sendWasmPackError(w, err)
		return
	}

	data, err := api.core.WasmConfigPack(deps_pack.Manifest)
	if SendIfError(w, err) {
		return
	}
	sum := sha256.Sum256(data)
	setWasmPackHeaders(w, `"config-`+hex.EncodeToString(sum[:])+`"`, schema.Name+".config.pack")
	http.ServeContent(w, r, "config.pack", time.Time{}, bytes.NewReader(data))
}
