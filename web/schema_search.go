package web

import (
	"blockexchange/types"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

func (api *Api) SearchRecentSchemas(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	count, err := strconv.Atoi(vars["count"])
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	list, err := api.SchemaSearchRepo.FindRecent(count)
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

	schema, err := api.SchemaSearchRepo.FindByUsernameAndSchemaname(schema_name, user_name)

	if r.URL.Query().Get("download") == "true" {
		// increment downloads and ignore error
		api.SchemaRepo.IncrementDownloads(schema.ID)
	}

	Send(w, schema, err)
}

var schemaSearchHistogram = promauto.NewHistogram(prometheus.HistogramOpts{
	Name:    "bx_schema_search_hist",
	Help:    "Histogram for the schema render time",
	Buckets: []float64{0.01, 0.05, 0.1, 0.5, 1, 2, 5, 10},
})

func (api *Api) SearchSchema(w http.ResponseWriter, r *http.Request) {
	timer := prometheus.NewTimer(schemaSearchHistogram)
	defer timer.ObserveDuration()

	search := types.SchemaSearch{}
	err := json.NewDecoder(r.Body).Decode(&search)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	if search.UserName != nil && search.SchemaName != nil {
		// direct search by username and schema name
		schema, err := api.SchemaSearchRepo.FindByUsernameAndSchemaname(*search.SchemaName, *search.UserName)
		Send(w, []types.SchemaSearchResult{*schema}, err)
		return
	}

	if search.UserID != nil {
		// search by user_id
		list, err := api.SchemaSearchRepo.FindByUserID(int64(*search.UserID))
		Send(w, list, err)
		return
	}

	if search.UserName != nil {
		// search by user_name
		list, err := api.SchemaSearchRepo.FindByUsername(*search.UserName)
		Send(w, list, err)
		return
	}

	if search.Keywords != nil {
		// search by keywords
		list, err := api.SchemaSearchRepo.FindByKeywords(*search.Keywords)
		Send(w, list, err)
		return
	}

	if search.SchemaID != nil {
		// search by schema id
		list, err := api.SchemaSearchRepo.FindBySchemaID(int64(*search.SchemaID))
		Send(w, list, err)
		return
	}

	if search.TagID != nil {
		// search by tag id
		list, err := api.SchemaSearchRepo.FindByTagID(int64(*search.TagID))
		Send(w, list, err)
		return
	}

}
