package api

import "net/http"

func (api *Api) Healthcheck(w http.ResponseWriter, r *http.Request) {
	if !api.Running.Load() {
		w.WriteHeader(http.StatusBadGateway)
		w.Write([]byte("Server is shutting down"))
		return
	}

	_, err := api.TagRepo.GetAll()
	Send(w, true, err)
}
