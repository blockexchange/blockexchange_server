package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

// SetupLegacyRedirects permanently redirects the old history-mode frontend urls
// (for example "/schema/{username}/{name}") to their hash-routed equivalent
// ("/#/schema/{username}/{name}") so existing links and bookmarks keep working.
//
// The list is frozen: new frontend routes are hash-based and need no entry here.
func (api *Api) SetupLegacyRedirects(r *mux.Router) {
	legacyUrls := []string{
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
		"/tags",
	}

	for _, u := range legacyUrls {
		r.HandleFunc(u, legacyRedirect)
	}
}

func legacyRedirect(w http.ResponseWriter, r *http.Request) {
	target := "/#" + r.URL.EscapedPath()
	if r.URL.RawQuery != "" {
		target += "?" + r.URL.RawQuery
	}
	http.Redirect(w, r, target, http.StatusMovedPermanently)
}
