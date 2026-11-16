package api

import (
	"blockexchange/public"
	"blockexchange/types"
	"net/http"
	"os"

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
		"/mods",
		"/schema/{username}",
		"/schema/{username}/{name}",
		"/user/{username}",
		"/collections/{username}",
		"/collections/{username}/{collection_name}",
		"/users",
		"/search",
		"/register",
		"/profile",
	}

	serveIndex := func(w http.ResponseWriter, r *http.Request) {
		var data []byte
		var err error
		if cfg.WebDev {
			// live mode: read from disk
			data, err = os.ReadFile("public/index.html")
		} else {
			data, err = public.Webapp.ReadFile("index.html")
		}
		if err != nil {
			SendError(w, 500, err.Error())
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(data)
	}

	for _, u := range indexUrls {
		r.HandleFunc(u, serveIndex)
	}
}
