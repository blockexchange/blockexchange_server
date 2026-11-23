package api

import (
	"blockexchange/types"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// returns all dependencies of a schema, schema owner or admin only
func (api *Api) GetSchemaModDependencies(w http.ResponseWriter, r *http.Request, ctx *SecureContext) {
	schema_uid := mux.Vars(r)["schema_uid"]

	if api.loadOwnedSchema(w, ctx, schema_uid, true) == nil {
		return
	}

	list, err := api.SchemaModDependencyRepo.GetBySchemaUID(schema_uid)
	Send(w, list, err)
}

// decodes and validates a dependency from the request body
func parseSchemaModDependency(w http.ResponseWriter, r *http.Request) *types.SchemaModDependency {
	dep := &types.SchemaModDependency{}
	err := json.NewDecoder(r.Body).Decode(dep)
	if err != nil {
		SendError(w, http.StatusBadRequest, err.Error())
		return nil
	}

	err = dep.NormalizeAndValidate()
	if err != nil {
		SendError(w, http.StatusBadRequest, err.Error())
		return nil
	}

	return dep
}

// loads the dependency from the url and checks that it belongs to the (owned) schema
func (api *Api) loadOwnedSchemaModDependency(w http.ResponseWriter, r *http.Request, ctx *SecureContext) *types.SchemaModDependency {
	vars := mux.Vars(r)
	schema_uid := vars["schema_uid"]

	if api.loadOwnedSchema(w, ctx, schema_uid, true) == nil {
		return nil
	}

	dep, err := api.SchemaModDependencyRepo.GetByUID(vars["uid"])
	if SendIfError(w, err) {
		return nil
	}
	if dep == nil || dep.SchemaUID != schema_uid {
		SendError(w, http.StatusNotFound, "dependency not found")
		return nil
	}
	return dep
}

// creates a new dependency, schema owner or admin only
func (api *Api) CreateSchemaModDependency(w http.ResponseWriter, r *http.Request, ctx *SecureContext) {
	schema_uid := mux.Vars(r)["schema_uid"]

	if api.loadOwnedSchema(w, ctx, schema_uid, true) == nil {
		return
	}

	dep := parseSchemaModDependency(w, r)
	if dep == nil {
		return
	}
	dep.UID = ""
	dep.SchemaUID = schema_uid

	err := api.SchemaModDependencyRepo.Create(dep)
	Send(w, dep, err)
}

// updates an existing dependency, schema owner or admin only
func (api *Api) UpdateSchemaModDependency(w http.ResponseWriter, r *http.Request, ctx *SecureContext) {
	existing := api.loadOwnedSchemaModDependency(w, r, ctx)
	if existing == nil {
		return
	}

	dep := parseSchemaModDependency(w, r)
	if dep == nil {
		return
	}
	dep.UID = existing.UID
	dep.SchemaUID = existing.SchemaUID

	err := api.SchemaModDependencyRepo.Update(dep)
	Send(w, dep, err)
}

// removes a dependency, schema owner or admin only
func (api *Api) DeleteSchemaModDependency(w http.ResponseWriter, r *http.Request, ctx *SecureContext) {
	dep := api.loadOwnedSchemaModDependency(w, r, ctx)
	if dep == nil {
		return
	}

	err := api.SchemaModDependencyRepo.Delete(dep.UID)
	if SendIfError(w, err) {
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
