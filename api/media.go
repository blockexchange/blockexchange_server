package api

import (
	"blockexchange/types"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// nodedef

func (api Api) CreateOrUpdateNodedefs(w http.ResponseWriter, r *http.Request, ctx *SecureContext) {
	if !ctx.CheckPermission(w, types.JWTPermissionMedia) {
		return
	}

	nodedefs := []*types.Nodedefinition{}
	err := json.NewDecoder(r.Body).Decode(&nodedefs)
	if err != nil {
		SendError(w, 500, fmt.Sprintf("json error: %s", err.Error()))
		return
	}

	for _, nd := range nodedefs {
		existing_nd, err := api.MediaRepo.GetNodedefinitionByName(nd.Name)
		if err != nil {
			SendError(w, 500, fmt.Sprintf("db error: %s", err.Error()))
			return
		}

		if existing_nd == nil {
			// insert
			err = api.MediaRepo.CreateNodedefinition(nd)
		} else {
			// update
			err = api.MediaRepo.UpdateNodedefinition(nd)
		}
		if err != nil {
			SendError(w, 500, fmt.Sprintf("update/insert error: %s", err.Error()))
			return
		}
	}

	Send(w, true, nil)
}

func (api Api) GetNodedefinition(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	m, err := api.MediaRepo.GetNodedefinitionByName(vars["nodename"])
	Send(w, m, err)
}

func (api Api) GetNodedefinitions(w http.ResponseWriter, r *http.Request) {
	m, err := api.MediaRepo.GetNodedefinitions()
	Send(w, m, err)
}

func (api Api) DeleteNodedefinition(w http.ResponseWriter, r *http.Request, ctx *SecureContext) {
	if !ctx.CheckPermission(w, types.JWTPermissionMedia) {
		return
	}

	vars := mux.Vars(r)
	err := api.MediaRepo.RemoveNodedefinition(vars["nodename"])
	Send(w, true, err)
}

// mediafile

func (api Api) CreateOrUpdateMediafiles(w http.ResponseWriter, r *http.Request, ctx *SecureContext) {
	if !ctx.CheckPermission(w, types.JWTPermissionMedia) {
		return
	}

	mediafiles := []*types.Mediafile{}
	err := json.NewDecoder(r.Body).Decode(&mediafiles)
	if err != nil {
		SendError(w, 500, fmt.Sprintf("json error: %s", err.Error()))
		return
	}

	for _, mediafile := range mediafiles {
		exsiting_mf, err := api.MediaRepo.GetMediafileByName(mediafile.Name)
		if err != nil {
			SendError(w, 500, fmt.Sprintf("db error: %s", err.Error()))
			return
		}

		if exsiting_mf == nil {
			// insert
			err = api.MediaRepo.CreateMediafile(mediafile)
		} else {
			// update
			err = api.MediaRepo.UpdateMediafile(mediafile)
		}
		if err != nil {
			SendError(w, 500, fmt.Sprintf("update/insert error: %s", err.Error()))
			return
		}
	}

	Send(w, true, nil)
}

func (api Api) GetMediafile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	m, err := api.MediaRepo.GetMediafileByName(vars["name"])
	if err != nil {
		SendError(w, 500, err.Error())
	} else if m == nil {
		SendError(w, 404, "not found")
	} else if r.URL.Query().Get("raw") == "true" {
		w.Write(m.Data)
	} else {
		Send(w, m, err)
	}
}

func (api Api) DeleteMediafile(w http.ResponseWriter, r *http.Request, ctx *SecureContext) {
	if !ctx.CheckPermission(w, types.JWTPermissionMedia) {
		return
	}

	vars := mux.Vars(r)
	err := api.MediaRepo.RemoveMediafile(vars["name"])
	Send(w, true, err)
}
