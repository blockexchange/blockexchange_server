package api

import (
	"blockexchange/types"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func (api *Api) GetSchemaMods(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	schema_uid := vars["schema_uid"]

	list, err := api.SchemaModRepo.GetSchemaModsBySchemaUID(schema_uid)
	if SendIfError(w, err) {
		return
	}

	modlist := make([]string, len(list))
	for i, mod := range list {
		modlist[i] = mod.ModName
	}

	SendJson(w, modlist)
}

func (api *Api) GetSchemaModCount(w http.ResponseWriter, r *http.Request) {
	list, err := api.SchemaModRepo.GetSchemaModCount()
	Send(w, list, err)
}

func (api *Api) CreateSchemaMods(w http.ResponseWriter, r *http.Request, ctx *SecureContext) {
	vars := mux.Vars(r)
	schema_uid := vars["schema_uid"]

	modlist := make([]string, 0)
	err := json.NewDecoder(r.Body).Decode(&modlist)
	if SendIfError(w, err) {
		return
	}

	if api.loadOwnedSchema(w, ctx, schema_uid, false) == nil {
		return
	}

	current_mod_list, err := api.SchemaModRepo.GetSchemaModsBySchemaUID(schema_uid)
	if SendIfError(w, err) {
		return
	}

	// collect existing mod names
	current_mod_map := make(map[string]bool)
	for _, mod_name := range current_mod_list {
		current_mod_map[mod_name.ModName] = true
	}

	for _, mod_name := range modlist {
		if current_mod_map[mod_name] {
			// name already exists, skip
			continue
		}
		err = api.SchemaModRepo.CreateSchemaMod(&types.SchemaMod{
			ModName:   mod_name,
			SchemaUID: schema_uid,
		})
		if SendIfError(w, err) {
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

func (api *Api) UpdateSchemaMods(w http.ResponseWriter, r *http.Request, ctx *SecureContext) {
	vars := mux.Vars(r)
	schema_uid := vars["schema_uid"]

	// security check
	if api.loadOwnedSchema(w, ctx, schema_uid, true) == nil {
		return
	}

	// extract modnames from schemaparts
	modnames, err := api.core.ExtractModnames(schema_uid)
	if SendIfError(w, err) {
		return
	}

	// remove old modnames
	err = api.Repositories.SchemaModRepo.RemoveSchemaMods(schema_uid)
	if SendIfError(w, err) {
		return
	}

	// add new modnames
	for _, modname := range modnames {
		err = api.Repositories.SchemaModRepo.CreateSchemaMod(&types.SchemaMod{
			SchemaUID: schema_uid,
			ModName:   modname,
		})
		if SendIfError(w, err) {
			return
		}
	}

	Send(w, modnames, nil)
}
