package api

import "net/http"

func (api Api) GetColorMapping(w http.ResponseWriter, r *http.Request) {
	SendJson(w, api.core.GetColorMapping())
}
