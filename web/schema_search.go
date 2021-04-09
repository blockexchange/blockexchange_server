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
		//TODO
	}

	if search.UserID != nil {
		// search by user_id
		//TODO
	}

	if search.UserName != nil {
		// search by user_name
		//TODO
	}

	if search.Keywords != nil {
		// search by keywords
		list, err := api.SchemaSearchRepo.FindByKeywords(*search.Keywords)
		Send(w, list, err)
		return
	}

	if search.SchemaID != nil {
		// search by schema id
		//TODO
	}

}
