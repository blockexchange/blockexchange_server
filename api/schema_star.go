package api

import (
	"blockexchange/types"
	"net/http"

	"github.com/gorilla/mux"
)

func (api *Api) CountSchemaStars(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	schema_uid := vars["schema_uid"]
	count, err := api.SchemaStarRepo.CountBySchemaUID(schema_uid)
	Send(w, count, err)
}

func (api *Api) GetSchemaStar(w http.ResponseWriter, r *http.Request, ctx *SecureContext) {
	vars := mux.Vars(r)
	schema_uid := vars["schema_uid"]
	star, err := api.SchemaStarRepo.GetBySchemaAndUserID(schema_uid, ctx.Claims.UserUID)
	Send(w, star, err)
}

func (api *Api) StarSchema(w http.ResponseWriter, r *http.Request, ctx *SecureContext) {
	vars := mux.Vars(r)
	schema_uid := vars["schema_uid"]

	err := api.SchemaStarRepo.Create(&types.SchemaStar{SchemaUID: schema_uid, UserUID: ctx.Claims.UserUID})
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	err = api.SchemaRepo.CalculateStats(schema_uid)
	Send(w, true, err)
}

func (api *Api) UnStarSchema(w http.ResponseWriter, r *http.Request, ctx *SecureContext) {
	vars := mux.Vars(r)
	schema_uid := vars["schema_uid"]

	err := api.SchemaStarRepo.Delete(&types.SchemaStar{SchemaUID: schema_uid, UserUID: ctx.Claims.UserUID})
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	err = api.SchemaRepo.CalculateStats(schema_uid)
	Send(w, true, err)
}
