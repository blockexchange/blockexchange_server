package web

import (
	"blockexchange/types"
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

	var has_user_id bool = true
	user_id, err := strconv.Atoi(r.URL.Query().Get("user_id"))
	if err != nil {
		// user not specified
		has_user_id = false
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

	response := &types.SchemaStarResponse{}

	// count the stars
	count, err := api.SchemaStarRepo.CountBySchemaID(int64(schema_id))
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}
	response.Count = count

	if has_user_id {
		// also check if the user has starred the schema already
		star, err := api.SchemaStarRepo.GetBySchemaAndUserID(int64(schema_id), int64(user_id))
		if err != nil {
			SendError(w, 500, err.Error())
			return
		}
		response.Starred = star != nil
	}

	SendJson(w, response)
}
