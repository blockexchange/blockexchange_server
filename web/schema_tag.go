package web

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (api Api) CreateSchemaTag(w http.ResponseWriter, r *http.Request, ctx *SecureContext) {
	vars := mux.Vars(r)
	schema_id, err := strconv.Atoi(vars["schema_id"])
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	tag_id, err := strconv.Atoi(vars["tag_id"])
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	schema, err := api.SchemaRepo.GetSchemaById(int64(schema_id))
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	if schema == nil {
		SendError(w, 404, "No schema found")
		return
	}

	if schema.UserID != ctx.Token.UserID {
		SendError(w, 403, "you are not the owner of the schema")
		return
	}

	err = api.SchemaTagRepo.Create(int64(schema_id), int64(tag_id))
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (api Api) DeleteSchemaTag(w http.ResponseWriter, r *http.Request, ctx *SecureContext) {
	vars := mux.Vars(r)
	schema_id, err := strconv.Atoi(vars["schema_id"])
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	tag_id, err := strconv.Atoi(vars["tag_id"])
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	schema, err := api.SchemaRepo.GetSchemaById(int64(schema_id))
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	if schema == nil {
		SendError(w, 404, "No schema found")
		return
	}

	if schema.UserID != ctx.Token.UserID {
		SendError(w, 403, "you are not the owner of the schema")
		return
	}

	err = api.SchemaTagRepo.Delete(int64(schema_id), int64(tag_id))
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (api Api) GetSchemaTags(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	schema_id, err := strconv.Atoi(vars["schema_id"])
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	list, err := api.SchemaTagRepo.GetBySchemaID(int64(schema_id))
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	SendJson(w, list)
}
