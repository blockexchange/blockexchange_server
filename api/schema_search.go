package api

import (
	"blockexchange/types"
	"encoding/json"
	"net/http"
)

func (api *Api) CountSchema(w http.ResponseWriter, r *http.Request) {
	search := &types.SchemaSearchRequest{}
	err := json.NewDecoder(r.Body).Decode(search)
	if SendIfError(w, err) {
		return
	}

	c, err := api.SchemaSearchRepo.Count(search)
	Send(w, c, err)
}

func (api *Api) SearchSchema(w http.ResponseWriter, r *http.Request) {
	search := &types.SchemaSearchRequest{}
	err := json.NewDecoder(r.Body).Decode(search)
	if SendIfError(w, err) {
		return
	}

	applySearchDefaults(search)
	list, err := api.SchemaSearchRepo.Search(search)
	if SendIfError(w, err) {
		return
	}

	Send(w, list, err)
}

// applies sane limit/offset defaults to a search request
func applySearchDefaults(search *types.SchemaSearchRequest) {
	if search.Limit == nil || *search.Limit > 100 || *search.Limit <= 0 {
		l := 100
		search.Limit = &l
	}
	if search.Offset == nil || *search.Offset > 10000 || *search.Offset < 0 {
		o := 0
		search.Offset = &o
	}
}

// searches the schematics starred by the logged-in user
func (api *Api) SearchStarredSchema(w http.ResponseWriter, r *http.Request, ctx *SecureContext) {
	search := &types.SchemaSearchRequest{}
	err := json.NewDecoder(r.Body).Decode(search)
	if SendIfError(w, err) {
		return
	}

	search.StarredByUserUID = &ctx.Claims.UserUID
	applySearchDefaults(search)
	list, err := api.SchemaSearchRepo.Search(search)
	Send(w, list, err)
}

// counts the schematics starred by the logged-in user
func (api *Api) CountStarredSchema(w http.ResponseWriter, r *http.Request, ctx *SecureContext) {
	search := &types.SchemaSearchRequest{}
	err := json.NewDecoder(r.Body).Decode(search)
	if SendIfError(w, err) {
		return
	}

	search.StarredByUserUID = &ctx.Claims.UserUID
	c, err := api.SchemaSearchRepo.Count(search)
	Send(w, c, err)
}
