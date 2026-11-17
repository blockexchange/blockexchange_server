package api

import (
	"blockexchange/schematic/bx"
	"blockexchange/schematic/worldedit"
	"blockexchange/types"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

const MAX_SIZE = 100 * 1000 * 1000

func (api *Api) ImportSchematic(w http.ResponseWriter, r *http.Request, ctx *SecureContext) {
	vars := mux.Vars(r)
	filename := vars["filename"]

	// abort reading as soon as the limit is exceeded instead of buffering the whole body
	data, err := io.ReadAll(http.MaxBytesReader(w, r.Body, MAX_SIZE))
	if err != nil {
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) {
			SendError(w, http.StatusRequestEntityTooLarge, fmt.Sprintf("filesize larger than %d bytes", MAX_SIZE))
			return
		}
		SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	var schema *types.Schema
	if strings.HasSuffix(filename, ".we") {
		schemaname := strings.TrimSuffix(filename, ".we")
		schema, err = api.core.ImportWE(data, ctx.Claims.Username, schemaname)
	} else if strings.HasSuffix(filename, ".zip") {
		schema, err = api.core.ImportBX(data, ctx.Claims.Username)
	} else {
		SendError(w, http.StatusBadRequest, "unrecognized file extension")
		return
	}

	if errors.Is(err, worldedit.ErrInvalidSchematic) || errors.Is(err, bx.ErrInvalidBX) {
		SendError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err != nil {
		SendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	schema, err = api.core.PostImport(schema)
	Send(w, schema, err)
}
