package api

import (
	"blockexchange/metrics"
	"blockexchange/parser"
	"blockexchange/types"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

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
	metrics.SchemaPartsUploaded.Inc()

	SendJson(w, schemapart)
}

// PatchSchemaPart applies the nodes selected by the patchfilter to the stored schemapart
func (api *Api) PatchSchemaPart(w http.ResponseWriter, r *http.Request, ctx *SecureContext) {
	patch := &types.PatchedSchemapart{}
	err := json.NewDecoder(r.Body).Decode(patch)
	if SendIfError(w, err) {
		return
	}

	schema := api.loadOwnedSchema(w, ctx, patch.SchemaUID, false)
	if schema == nil {
		return
	}

	schemapart, err := api.SchemaPartRepo.GetBySchemaUIDAndOffset(patch.SchemaUID, patch.OffsetX, patch.OffsetY, patch.OffsetZ)
	if SendIfError(w, err) {
		return
	}
	if schemapart == nil {
		SendError(w, http.StatusNotFound, "schemapart not found")
		return
	}

	parsed, err := parser.ParseSchemaPart(schemapart)
	if err != nil {
		SendError(w, 500, fmt.Sprintf("parse stored schemapart error: %s", err))
		return
	}

	parsed_patch, err := parser.ParseSchemaPart(patch.SchemaPart)
	if err != nil {
		SendError(w, http.StatusBadRequest, fmt.Sprintf("parse patch error: %s", err))
		return
	}

	err = parsed.Patch(parsed_patch, patch.PatchFilter)
	if err != nil {
		SendError(w, http.StatusBadRequest, fmt.Sprintf("patch error: %s", err))
		return
	}

	patched, err := parsed.Convert()
	if err != nil {
		SendError(w, 500, fmt.Sprintf("convert error: %s", err))
		return
	}

	mtime := time.Now().UnixMilli()

	// update schema part
	patched.SchemaUID = schemapart.SchemaUID
	patched.Mtime = mtime
	err = api.SchemaPartRepo.CreateOrUpdateSchemaPart(patched)
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

	metrics.SchemaPartsUploaded.Inc()

	SendJson(w, patched)
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

	metrics.SchemaPartsDownloaded.Inc()

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
		metrics.SchemaPartsDownloaded.Add(float64(len(schemaparts)))
		Send(w, schemaparts, err)
	}
}

func (api *Api) GetNextSchemaPart(w http.ResponseWriter, r *http.Request) {
	schema_uid, x, y, z, err := extractSchemaPartVars(r)
	if SendIfError(w, err) {
		return
	}

	metrics.SchemaPartsDownloaded.Inc()

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

	metrics.SchemaPartsDownloaded.Inc()

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

	metrics.SchemaPartsDownloaded.Inc()

	schemapart, err := api.SchemaPartRepo.GetFirstBySchemaUID(schema_uid)
	if err == nil && schemapart == nil {
		w.WriteHeader(http.StatusNoContent)
	} else {
		Send(w, schemapart, err)
	}
}
