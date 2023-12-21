package api

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (api *Api) CountSchemaStars(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	schema_id, err := strconv.ParseInt(vars["schema_id"], 10, 64)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}
	count, err := api.SchemaStarRepo.CountBySchemaID(schema_id)
	Send(w, count, err)
}

func (api *Api) GetSchemaStar(w http.ResponseWriter, r *http.Request, ctx *SecureContext) {
	vars := mux.Vars(r)
	schema_id, err := strconv.ParseInt(vars["schema_id"], 10, 64)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}
	star, err := api.SchemaStarRepo.GetBySchemaAndUserID(schema_id, ctx.Claims.UserID)
	Send(w, star, err)
}

func (api *Api) StarSchema(w http.ResponseWriter, r *http.Request, ctx *SecureContext) {
	vars := mux.Vars(r)
	schema_id, err := strconv.ParseInt(vars["schema_id"], 10, 64)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	err = api.SchemaStarRepo.Create(schema_id, ctx.Claims.UserID)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	err = api.SchemaRepo.CalculateStats(schema_id)
	Send(w, true, err)
}

func (api *Api) UnStarSchema(w http.ResponseWriter, r *http.Request, ctx *SecureContext) {
	vars := mux.Vars(r)
	schema_id, err := strconv.ParseInt(vars["schema_id"], 10, 64)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	err = api.SchemaStarRepo.Delete(schema_id, ctx.Claims.UserID)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	err = api.SchemaRepo.CalculateStats(schema_id)
	Send(w, true, err)
}
