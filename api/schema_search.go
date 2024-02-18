package api

import (
	"blockexchange/types"
	"encoding/json"
	"net/http"
)

func (api *Api) CountSchema(w http.ResponseWriter, r *http.Request) {
	search := &types.SchemaSearchRequest{}
	err := json.NewDecoder(r.Body).Decode(search)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	c, err := api.SchemaSearchRepo.Count(search)
	Send(w, c, err)
}

func (api *Api) SearchSchema(w http.ResponseWriter, r *http.Request) {
	search := &types.SchemaSearchRequest{}
	err := json.NewDecoder(r.Body).Decode(search)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	// apply sane defaults
	if search.Limit == nil || *search.Limit > 100 || *search.Limit <= 0 {
		l := 100
		search.Limit = &l
	}
	if search.Offset == nil || *search.Offset > 10000 || *search.Offset < 0 {
		o := 0
		search.Offset = &o
	}

	list, err := api.SchemaSearchRepo.Search(search)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	Send(w, list, err)
}
