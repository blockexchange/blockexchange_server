package web

import "net/http"

func (api *Api) Index(w http.ResponseWriter, r *http.Request) {
	api.te.Execute("pages/index.html", w, r, "./", nil)
}
