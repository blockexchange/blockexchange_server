package api

import (
	"blockexchange/types"
	"net/http"

	"github.com/gorilla/mux"
)

func (api *Api) SearchSchemaByNameAndUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	schema_name := vars["schema_name"]
	user_name := vars["user_name"]

	search := &types.SchemaSearchRequest{
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
