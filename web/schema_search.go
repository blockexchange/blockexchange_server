package web

import (
	"blockexchange/types"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func (api *Api) SearchRecentSchemas(w http.ResponseWriter, r *http.Request) {
	limit, offset := GetLimitOffset(r, 20)

	complete := true
	search := &types.SchemaSearch{
		Complete: &complete,
	}
	list, err := api.SchemaSearchRepo.Search(search, limit, offset)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	SendJson(w, list)
}

func (api *Api) SearchSchemaByNameAndUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	schema_name := vars["schema_name"]
	user_name := vars["user_name"]

	search := &types.SchemaSearch{
		UserName:   &user_name,
		SchemaName: &schema_name,
	}
	list, err := api.SchemaSearchRepo.Search(search, 1, 0)
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

func (api *Api) SearchSchema(w http.ResponseWriter, r *http.Request) {
	limit, offset := GetLimitOffset(r, 20)

	search := &types.SchemaSearch{}
	err := json.NewDecoder(r.Body).Decode(search)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	list, err := api.SchemaSearchRepo.Search(search, limit, offset)
	Send(w, list, err)
}
