package api

import (
	"blockexchange/public"
	"blockexchange/types"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func (api *Api) SetupSPARoutes(r *mux.Router, cfg *types.Config) {
	// index.html or urls that point to it
	indexUrls := []string{
		"/",
		"/index.html",
		"/login",
		"/user/{username}",
		"/schema/{username}",
		"/users",
		"/search",
		"/register",
		"/profile",
	}
	for _, u := range indexUrls {
		r.HandleFunc(u, func(w http.ResponseWriter, r *http.Request) {
			public.RenderIndex(w, r, map[string]string{
				"og:site_name": "Blockexchange",
			})
		})
	}

	r.HandleFunc("/schema/{username}/{name}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		username := vars["username"]
		schema, err := api.SchemaRepo.GetSchemaByUsernameAndName(username, vars["name"])
		if err != nil {
			SendError(w, 500, err.Error())
			return
		}
		if schema == nil {
			SendError(w, 404, "not found")
			return
		}

		meta := map[string]string{
			"og:site_name":   "Blockexchange",
			"og:description": schema.Description,
			"og:title":       fmt.Sprintf("'%s' by %s", schema.Name, username),
			"og:type":        "Schematic",
			"og:url":         fmt.Sprintf("%s/schema/%s/%s", api.cfg.BaseURL, username, schema.Name),
			"og:image":       fmt.Sprintf("%s/api/schema/%d/screenshot", api.cfg.BaseURL, schema.ID),
		}

		public.RenderIndex(w, r, meta)
	})
}
