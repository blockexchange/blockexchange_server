package api

import (
	"blockexchange/types"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func (api Api) GetSchemaPull(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	schema_uid := vars["schema_uid"]

	list, err := api.Repositories.SchemaPullRepo.GetBySchemaUID(schema_uid)
	Send(w, list, err)
}

func (api Api) CreateSchemaPull(w http.ResponseWriter, r *http.Request, ctx *SecureContext) {
	vars := mux.Vars(r)
	schema_uid := vars["schema_uid"]

	schema, err := api.SchemaRepo.GetSchemaByUID(schema_uid)
	if err != nil {
		SendError(w, 500, fmt.Sprintf("getby schemaUID error: %s", err))
		return
	}

	if schema == nil {
		SendError(w, 500, "no schema found")
		return
	}

	if schema.UserUID != ctx.Claims.UserUID {
		SendError(w, 403, "you are not the owner of the schema")
		return
	}

	sp := &types.SchematicPull{}
	err = json.NewDecoder(r.Body).Decode(sp)
	if err != nil {
		SendError(w, 500, fmt.Sprintf("json error: %v", err))
		return
	}
	sp.SchemaUID = schema.UID

	err = api.Repositories.SchemaPullRepo.Create(sp)
	Send(w, sp, err)
}

func (api Api) GetSchemaPullClients(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	schema_uid := vars["schema_uid"]

	list, err := api.Repositories.SchemaPullClientRepo.GetBySchemaUID(schema_uid)
	Send(w, list, err)
}
