package web

import (
	"net/http"
)

func (api *Api) GetTags(w http.ResponseWriter, req *http.Request) {
	list, err := api.TagRepo.GetAll()
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}
	SendJson(w, list)
}
