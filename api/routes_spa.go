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
		"/import",
		"/mod",
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
			"og:description": schema.ShortDescription,
			"og:title":       fmt.Sprintf("'%s' by %s", schema.Name, username),
			"og:type":        "Schematic",
			"og:url":         fmt.Sprintf("%s/schema/%s/%s", api.cfg.BaseURL, username, schema.Name),
			"og:image":       fmt.Sprintf("%s/api/schema/%d/screenshot", api.cfg.BaseURL, schema.ID),
		}

		public.RenderIndex(w, r, meta)
	})

	r.HandleFunc("/user/{username}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		username := vars["username"]
		user, err := api.UserRepo.GetUserByName(username)
		if err != nil {
			SendError(w, 500, err.Error())
			return
		}

		schematics, err := api.SchemaSearchRepo.Count(&types.SchemaSearchRequest{
			UserID: user.ID,
		})
		if err != nil {
			SendError(w, 500, err.Error())
			return
		}

		stars, err := api.SchemaStarRepo.CountByUserID(*user.ID)
		if err != nil {
			SendError(w, 500, err.Error())
			return
		}

		meta := map[string]string{
			"og:site_name":   "Blockexchange",
			"og:title":       fmt.Sprintf("User '%s'", username),
			"og:description": fmt.Sprintf("Schematics: %d, Stars: %d ★", schematics, stars),
			"og:url":         fmt.Sprintf("%s/user/%s", api.cfg.BaseURL, username),
			"og:image":       user.AvatarURL,
		}

		public.RenderIndex(w, r, meta)
	})
}