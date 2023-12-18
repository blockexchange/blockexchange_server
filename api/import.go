package api

import (
	"blockexchange/types"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

const MAX_SIZE = 30 * 1000 * 1000

func (api *Api) ImportSchematic(w http.ResponseWriter, r *http.Request, ctx *SecureContext) {
	vars := mux.Vars(r)
	filename := vars["filename"]

	data, err := io.ReadAll(r.Body)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}
	if len(data) > MAX_SIZE {
		SendError(w, 500, fmt.Sprintf("filesize larger that %d", MAX_SIZE))
		return
	}

	var schema *types.Schema
	if strings.HasSuffix(filename, ".we") {
		schemaname := strings.TrimSuffix(filename, ".we")
		schema, err = api.core.ImportWE(data, ctx.Claims.Username, schemaname)
	} else if strings.HasSuffix(filename, ".zip") {
		schema, err = api.core.ImportBX(data, ctx.Claims.Username)
	} else {
		err = errors.New("unrecognized file extension")
	}

	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	schema, err = api.core.PostImport(schema)
	Send(w, schema, err)
}
