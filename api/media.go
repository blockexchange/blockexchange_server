package api

import (
	"blockexchange/db"
	"blockexchange/types"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

const maxMediaPageSize = 500

// max number of names in a batch lookup
const maxMediaBatchSize = 1000

// parses a json array of names from the request body
func parseNameBatch(w http.ResponseWriter, r *http.Request) ([]string, bool) {
	names := []string{}
	err := json.NewDecoder(r.Body).Decode(&names)
	if err != nil {
		SendError(w, 400, fmt.Sprintf("json error: %s", err.Error()))
		return nil, false
	}
	if len(names) > maxMediaBatchSize {
		SendError(w, 400, fmt.Sprintf("too many names: %d, max: %d", len(names), maxMediaBatchSize))
		return nil, false
	}
	return names, true
}

// parses the optional "q", "limit" and "offset" query params, no limit means all entries
func parseMediaQuery(r *http.Request) db.MediaQuery {
	q := db.MediaQuery{
		Query: r.URL.Query().Get("q"),
		Limit: -1,
	}
	if limit, err := strconv.Atoi(r.URL.Query().Get("limit")); err == nil && limit > 0 {
		q.Limit = min(limit, maxMediaPageSize)
	}
	if offset, err := strconv.Atoi(r.URL.Query().Get("offset")); err == nil && offset > 0 {
		q.Offset = offset
	}
	return q
}

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

// returns the node definitions for a json array of node names, unknown names are skipped
func (api Api) GetNodedefinitionsBatch(w http.ResponseWriter, r *http.Request) {
	names, ok := parseNameBatch(w, r)
	if !ok {
		return
	}
	list, err := api.MediaRepo.GetNodedefinitionsByNames(names)
	Send(w, list, err)
}

func (api Api) GetNodedefinitions(w http.ResponseWriter, r *http.Request) {
	m, err := api.MediaRepo.SearchNodedefinitions(parseMediaQuery(r))
	Send(w, m, err)
}

func (api Api) CountNodedefinitions(w http.ResponseWriter, r *http.Request) {
	c, err := api.MediaRepo.CountNodedefinitions(r.URL.Query().Get("q"))
	Send(w, c, err)
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

func (api Api) GetMediafileInfos(w http.ResponseWriter, r *http.Request) {
	m, err := api.MediaRepo.SearchMediafileInfos(parseMediaQuery(r))
	Send(w, m, err)
}

func (api Api) CountMediafiles(w http.ResponseWriter, r *http.Request) {
	c, err := api.MediaRepo.CountMediafiles(r.URL.Query().Get("q"))
	Send(w, c, err)
}

// returns the media files for a json array of file names, unknown names are skipped
func (api Api) GetMediafilesBatch(w http.ResponseWriter, r *http.Request) {
	names, ok := parseNameBatch(w, r)
	if !ok {
		return
	}
	list, err := api.MediaRepo.GetMediafilesByNames(names)
	Send(w, list, err)
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
