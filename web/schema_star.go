package web

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (api Api) CreateSchemaStar(w http.ResponseWriter, r *http.Request, ctx *SecureContext) {
	vars := mux.Vars(r)
	schema_id, err := strconv.Atoi(vars["schema_id"])
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

	err = api.SchemaStarRepo.Create(int64(schema_id), ctx.Token.UserID)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (api Api) DeleteSchemaStar(w http.ResponseWriter, r *http.Request, ctx *SecureContext) {
	vars := mux.Vars(r)
	schema_id, err := strconv.Atoi(vars["schema_id"])
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

	err = api.SchemaStarRepo.Delete(int64(schema_id), ctx.Token.UserID)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (api Api) GetSchemaStars(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	schema_id, err := strconv.Atoi(vars["schema_id"])
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

	count, err := api.SchemaStarRepo.CountBySchemaID(int64(schema_id))
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("content-type", "text/plain")
	w.Write([]byte(strconv.Itoa(count)))
}
