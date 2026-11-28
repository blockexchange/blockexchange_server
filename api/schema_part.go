package api

import (
	"blockexchange/db"
	"blockexchange/metrics"
	"blockexchange/parser"
	"blockexchange/types"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

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

	// create or update the schema part with a unique mtime (updates the schema mtime too)
	err = api.SchemaPartRepo.UpdateWithNextMtime(schema.UID, func(repo *db.SchemaPartRepository, mtime int64) error {
		schemapart.Mtime = mtime
		return repo.CreateOrUpdateSchemaPart(schemapart)
	})
	if err != nil {
		SendError(w, 500, fmt.Sprintf("CreateOrUpdateSchemaPart error: %s", err))
		return
	}

	// increment stats
	metrics.SchemaPartsUploaded.Inc()

	SendJson(w, schemapart)
}

// statusError is an error with an http status code
type statusError struct {
	code int
	msg  string
}

func (e *statusError) Error() string {
	return e.msg
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

	parsed_patch, err := parser.ParseSchemaPart(patch.SchemaPart)
	if err != nil {
		SendError(w, http.StatusBadRequest, fmt.Sprintf("parse patch error: %s", err))
		return
	}

	// read, patch and write the stored part in one transaction: concurrent patches of the same schema are serialized
	var patched *types.SchemaPart
	err = api.SchemaPartRepo.UpdateWithNextMtime(schema.UID, func(repo *db.SchemaPartRepository, mtime int64) error {
		schemapart, err := repo.GetBySchemaUIDAndOffset(patch.SchemaUID, patch.OffsetX, patch.OffsetY, patch.OffsetZ)
		if err != nil {
			return err
		}
		if schemapart == nil {
			return &statusError{code: http.StatusNotFound, msg: "schemapart not found"}
		}

		parsed, err := parser.ParseSchemaPart(schemapart)
		if err != nil {
			return fmt.Errorf("parse stored schemapart error: %s", err)
		}

		err = parsed.Patch(parsed_patch, patch.PatchFilter)
		if err != nil {
			return &statusError{code: http.StatusBadRequest, msg: fmt.Sprintf("patch error: %s", err)}
		}

		patched, err = parsed.Convert()
		if err != nil {
			return fmt.Errorf("convert error: %s", err)
		}

		patched.SchemaUID = schemapart.SchemaUID
		patched.Mtime = mtime
		return repo.CreateOrUpdateSchemaPart(patched)
	})

	var status_err *statusError
	if errors.As(err, &status_err) {
		SendError(w, status_err.code, status_err.msg)
		return
	} else if err != nil {
		SendError(w, 500, err.Error())
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

	// remove the part and update the schema mtime
	err = api.SchemaPartRepo.UpdateWithNextMtime(schema_uid, func(repo *db.SchemaPartRepository, _ int64) error {
		return repo.RemoveBySchemaUIDAndOffset(schema_uid, x, y, z)
	})
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
