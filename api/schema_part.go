package api

import (
	"blockexchange/types"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var partsUploaded = promauto.NewCounter(prometheus.CounterOpts{
	Name: "bx_schemaparts_uploaded",
	Help: "The total number of uploaded schemaparts",
})

var partsDownloaded = promauto.NewCounter(prometheus.CounterOpts{
	Name: "bx_schemaparts_downloaded",
	Help: "The total number of downloaded schemaparts",
})

func (api *Api) CreateSchemaPart(w http.ResponseWriter, r *http.Request, ctx *SecureContext) {
	schemapart := &types.SchemaPart{}
	err := json.NewDecoder(r.Body).Decode(schemapart)
	if SendIfError(w, err) {
		return
	}

	schema := api.loadOwnedSchema(w, ctx, schemapart.SchemaUID, false)
	if schema == nil {
		return
	}

	mtime := time.Now().UnixMilli()

	// update schema part
	schemapart.Mtime = mtime
	err = api.SchemaPartRepo.CreateOrUpdateSchemaPart(schemapart)
	if err != nil {
		SendError(w, 500, fmt.Sprintf("CreateOrUpdateSchemaPart error: %s", err))
		return
	}

	// update schema mtime
	schema.Mtime = mtime
	err = api.SchemaRepo.UpdateSchema(schema)
	if err != nil {
		SendError(w, 500, fmt.Sprintf("UpdateSchema error: %s", err))
		return
	}

	// increment stats
	partsUploaded.Inc()

	SendJson(w, schemapart)
}

func (api *Api) DeleteSchemaPart(w http.ResponseWriter, r *http.Request, ctx *SecureContext) {
	schema_uid, x, y, z, err := extractSchemaPartVars(r)
	if SendIfError(w, err) {
		return
	}

	schema := api.loadOwnedSchema(w, ctx, schema_uid, false)
	if schema == nil {
		return
	}

	err = api.SchemaPartRepo.RemoveBySchemaUIDAndOffset(schema_uid, x, y, z)
	if SendIfError(w, err) {
		return
	}

	// update schema mtime
	schema.Mtime = time.Now().UnixMilli()
	err = api.SchemaRepo.UpdateSchema(schema)
	if SendIfError(w, err) {
		return
	}

	w.Write([]byte("true"))
}

func extractSchemaPartVars(r *http.Request) (string, int, int, int, error) {
	vars := mux.Vars(r)
	schema_uid := vars["schema_uid"]

	x, err := strconv.Atoi(vars["x"])
	if err != nil {
		return "", 0, 0, 0, err
	}

	y, err := strconv.Atoi(vars["y"])
	if err != nil {
		return "", 0, 0, 0, err
	}

	z, err := strconv.Atoi(vars["z"])
	if err != nil {
		return "", 0, 0, 0, err
	}

	return schema_uid, x, y, z, nil
}

func (api *Api) GetSchemaPart(w http.ResponseWriter, r *http.Request) {
	schema_uid, x, y, z, err := extractSchemaPartVars(r)
	if SendIfError(w, err) {
		return
	}

	partsDownloaded.Inc()

	schemapart, err := api.SchemaPartRepo.GetBySchemaUIDAndOffset(schema_uid, x, y, z)
	if err == nil && schemapart == nil {
		w.WriteHeader(http.StatusNoContent)
	} else {
		Send(w, schemapart, err)
	}
}

func (api *Api) GetSchemaPartChunk(w http.ResponseWriter, r *http.Request) {
	schema_uid, x, y, z, err := extractSchemaPartVars(r)
	if SendIfError(w, err) {
		return
	}

	// a chunk covers 4x4x4 parts, the range is inclusive so stop before the next chunk
	to_pos_offset := (16 * 4) - 1
	schemaparts, err := api.SchemaPartRepo.GetBySchemaUIDAndRange(schema_uid, x, y, z, x+to_pos_offset, y+to_pos_offset, z+to_pos_offset)
	if err == nil && schemaparts == nil {
		w.WriteHeader(http.StatusNoContent)
	} else {
		partsDownloaded.Add(float64(len(schemaparts)))
		Send(w, schemaparts, err)
	}
}

func (api *Api) GetNextSchemaPart(w http.ResponseWriter, r *http.Request) {
	schema_uid, x, y, z, err := extractSchemaPartVars(r)
	if SendIfError(w, err) {
		return
	}

	partsDownloaded.Inc()

	schemapart, err := api.SchemaPartRepo.GetNextBySchemaUIDAndOffset(schema_uid, x, y, z)
	if err == nil && schemapart == nil {
		w.WriteHeader(http.StatusNoContent)
	} else {
		Send(w, schemapart, err)
	}
}

func (api *Api) GetNextSchemaPartByMtime(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	schema_uid := vars["schema_uid"]

	mtime, err := strconv.Atoi(vars["mtime"])
	if SendIfError(w, err) {
		return
	}

	partsDownloaded.Inc()

	schemapart, err := api.SchemaPartRepo.GetNextBySchemaUIDAndMtime(schema_uid, int64(mtime))
	if err == nil && schemapart == nil {
		w.WriteHeader(http.StatusNoContent)
	} else {
		Send(w, schemapart, err)
	}
}

func (api *Api) CountNextSchemaPartByMtime(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	schema_uid := vars["schema_uid"]

	mtime, err := strconv.Atoi(vars["mtime"])
	if SendIfError(w, err) {
		return
	}

	count, err := api.SchemaPartRepo.CountNextBySchemaUIDAndMtime(schema_uid, int64(mtime))
	if SendIfError(w, err) {
		return
	}

	w.Write([]byte(fmt.Sprintf("%d", count)))
}

func (api *Api) GetFirstSchemaPart(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	schema_uid := vars["schema_uid"]

	partsDownloaded.Inc()

	schemapart, err := api.SchemaPartRepo.GetFirstBySchemaUID(schema_uid)
	if err == nil && schemapart == nil {
		w.WriteHeader(http.StatusNoContent)
	} else {
		Send(w, schemapart, err)
	}
}
