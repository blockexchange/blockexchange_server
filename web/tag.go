package web

import (
	"blockexchange/types"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (api *Api) GetTags(w http.ResponseWriter, req *http.Request) {
	list, err := api.TagRepo.GetAll()
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}
	SendJson(w, list)
}

func (api *Api) CreateTag(w http.ResponseWriter, r *http.Request, ctx *SecureContext) {
	if !ctx.CheckPermission(w, types.JWTPermissionAdmin) {
		return
	}

	tag := types.Tag{}
	err := json.NewDecoder(r.Body).Decode(&tag)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	err = api.TagRepo.Create(&tag)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	SendJson(w, tag)
}

func (api *Api) UpdateTag(w http.ResponseWriter, r *http.Request, ctx *SecureContext) {
	if !ctx.CheckPermission(w, types.JWTPermissionAdmin) {
		return
	}

	tag := types.Tag{}
	err := json.NewDecoder(r.Body).Decode(&tag)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	err = api.TagRepo.Update(&tag)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	SendJson(w, tag)
}

func (api *Api) DeleteTag(w http.ResponseWriter, r *http.Request, ctx *SecureContext) {
	if !ctx.CheckPermission(w, types.JWTPermissionAdmin) {
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	err = api.TagRepo.Delete(int64(id))
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}
