package web

import (
	"io/fs"
	"net/http"
)

func HandleAssets(assets fs.FS, cache bool) http.HandlerFunc {
	hs := http.FileServer(http.FS(assets))
	return func(w http.ResponseWriter, r *http.Request) {
		if cache {
			// cache all assets
			w.Header().Add("cache-control", "max-age=345600")
		}
		hs.ServeHTTP(w, r)
	}
}
