package api

import (
	"blockexchange/core"
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// max time for cloning all dependencies of a schematic
const wasmPackPrepareTimeout = 5 * time.Minute

// returns a zstd compressed tar archive with the game and mods of the schematic for the luanti wasm client
func (api *Api) GetWasmPack(w http.ResponseWriter, r *http.Request, ctx *SecureContext) {
	schema_uid := mux.Vars(r)["schema_uid"]

	schema := api.loadOwnedSchema(w, ctx, schema_uid, true)
	if schema == nil {
		return
	}

	deps, err := api.SchemaModDependencyRepo.GetBySchemaUID(schema_uid)
	if SendIfError(w, err) {
		return
	}
	if len(deps) == 0 {
		SendError(w, http.StatusNotFound, "the schematic has no mod dependencies")
		return
	}

	prepare_ctx, cancel := context.WithTimeout(r.Context(), wasmPackPrepareTimeout)
	defer cancel()
	pack, err := api.core.PrepareWasmPack(prepare_ctx, schema.Name, deps)
	switch {
	case errors.Is(err, core.ErrWasmPackInvalid):
		SendError(w, http.StatusBadRequest, err.Error())
		return
	case errors.Is(err, context.DeadlineExceeded):
		SendError(w, http.StatusGatewayTimeout, "cloning the dependencies took too long")
		return
	case err != nil:
		// clone failures (unreachable repository, unknown version)
		SendError(w, http.StatusBadGateway, err.Error())
		return
	}
	defer pack.Close()

	w.Header().Set("Content-Type", "application/zstd")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s.pack\"", schema.Name))
	err = pack.Write(w)
	if err != nil {
		// headers are already sent, just log
		logrus.WithError(err).WithField("schema_uid", schema_uid).Error("writing wasm pack failed")
	}
}
