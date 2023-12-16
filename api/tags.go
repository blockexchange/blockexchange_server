package api

import "net/http"

func (api Api) GetTags(w http.ResponseWriter, r *http.Request) {
	tags, err := api.TagRepo.GetAll()
	Send(w, tags, err)
}
