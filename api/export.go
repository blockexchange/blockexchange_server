package api

import (
	"blockexchange/core"
	"blockexchange/types"
	"blockexchange/worldedit"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (api *Api) ExportWorldeditSchema(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	var schemapart *types.SchemaPart
	it := func() (*types.SchemaPart, error) {
		var err error
		if schemapart == nil {
			schemapart, err = api.SchemaPartRepo.GetFirstBySchemaID(int64(id))
		} else {
			schemapart, err = api.SchemaPartRepo.GetNextBySchemaIDAndOffset(int64(id), schemapart.OffsetX, schemapart.OffsetY, schemapart.OffsetZ)
		}
		return schemapart, err
	}

	err = worldedit.Export(w, it)
	if err != nil {
		SendError(w, 500, err.Error())
	}

	err = api.SchemaRepo.IncrementDownloads(int64(id))
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}
}

func (api *Api) ExportBXSchema(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	schema, err := api.SchemaRepo.GetSchemaById(int64(id))
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}
	if schema == nil {
		SendError(w, 404, "not found")
		return
	}

	schemamods, err := api.SchemaModRepo.GetSchemaModsBySchemaID(int64(id))
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	var schemapart *types.SchemaPart
	it := func() (*types.SchemaPart, error) {
		var err error
		if schemapart == nil {
			schemapart, err = api.SchemaPartRepo.GetFirstBySchemaID(int64(id))
		} else {
			schemapart, err = api.SchemaPartRepo.GetNextBySchemaIDAndOffset(int64(id), schemapart.OffsetX, schemapart.OffsetY, schemapart.OffsetZ)
		}
		return schemapart, err
	}

	err = core.ExportBXSchema(w, schema, schemamods, it)
	if err != nil {
		SendError(w, 500, err.Error())
	}

	err = api.SchemaRepo.IncrementDownloads(int64(id))
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}
}
