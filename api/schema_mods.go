package api

import (
	"blockexchange/types"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (api *Api) GetSchemaMods(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["schema_id"])
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	list, err := api.SchemaModRepo.GetSchemaModsBySchemaID(int64(id))
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	modlist := make([]string, len(list))
	for i, mod := range list {
		modlist[i] = mod.ModName
	}

	SendJson(w, modlist)
}

func (api *Api) CreateSchemaMods(w http.ResponseWriter, r *http.Request, ctx *SecureContext) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["schema_id"])
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	modlist := make([]string, 0)
	err = json.NewDecoder(r.Body).Decode(&modlist)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	schema, err := api.SchemaRepo.GetSchemaById(int64(id))
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	if schema.UserID != ctx.Claims.UserID {
		SendError(w, 403, "you are not the owner of the schema")
		return
	}

	current_mod_list, err := api.SchemaModRepo.GetSchemaModsBySchemaID(int64(id))
	if err != nil {
		SendError(w, 500, err.Error())
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
			ModName:  mod_name,
			SchemaID: int64(id),
		})
		if err != nil {
			SendError(w, 500, err.Error())
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}
