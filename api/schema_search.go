package api

import (
	"blockexchange/types"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func (api *Api) SearchSchemaByNameAndUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	schema_name := vars["schema_name"]
	user_name := vars["user_name"]
	limit := 1
	offset := 0

	search := &types.SchemaSearchRequest{
		UserName:   &user_name,
		SchemaName: &schema_name,
		Limit:      &limit,
		Offset:     &offset,
	}
	list, err := api.SchemaSearchRepo.Search(search)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	if len(list) == 0 {
		SendError(w, 404, "not found")
		return
	}

	schema := list[0]
	if r.URL.Query().Get("download") == "true" {
		// increment downloads and ignore error
		api.SchemaRepo.IncrementDownloads(schema.ID)
	}

	Send(w, schema, err)
}

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
	Send(w, list, err)
}
