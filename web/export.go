package web

import (
	"blockexchange/core"
	"blockexchange/types"
	"net/http"

	"github.com/gorilla/mux"
)

func (api *Api) ExportWorldeditSchema(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	schema, err := api.SchemaRepo.GetSchemaByUsernameAndName(vars["username"], vars["schemaname"])
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	if schema == nil {
		SendError(w, 404, "Not found")
		return
	}

	var schemapart *types.SchemaPart
	it := func() (*types.SchemaPart, error) {
		var err error
		if schemapart == nil {
			schemapart, err = api.SchemaPartRepo.GetFirstBySchemaID(schema.ID)
		} else {
			schemapart, err = api.SchemaPartRepo.GetNextBySchemaIDAndOffset(schema.ID, schemapart.OffsetX, schemapart.OffsetY, schemapart.OffsetZ)
		}
		return schemapart, err
	}

	err = core.ExportWorldeditSchema(w, it)
	if err != nil {
		SendError(w, 500, err.Error())
	}
}

func (api *Api) ExportBXSchema(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	schema, err := api.SchemaRepo.GetSchemaByUsernameAndName(vars["username"], vars["schemaname"])
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	if schema == nil {
		SendError(w, 404, "Not found")
		return
	}

	schemamods, err := api.SchemaModRepo.GetSchemaModsBySchemaID(schema.ID)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	var schemapart *types.SchemaPart
	it := func() (*types.SchemaPart, error) {
		var err error
		if schemapart == nil {
			schemapart, err = api.SchemaPartRepo.GetFirstBySchemaID(schema.ID)
		} else {
			schemapart, err = api.SchemaPartRepo.GetNextBySchemaIDAndOffset(schema.ID, schemapart.OffsetX, schemapart.OffsetY, schemapart.OffsetZ)
		}
		return schemapart, err
	}

	err = core.ExportBXSchema(w, schema, schemamods, it)
	if err != nil {
		SendError(w, 500, err.Error())
	}
}
