package api

import (
	"blockexchange/core"
	"net/http"
)

func (api *Api) GetCDBVcsPackages(w http.ResponseWriter, r *http.Request) {
	data, err := core.CDBVcsPackages.Get()
	if err != nil {
		SendError(w, http.StatusBadGateway, err.Error())
		return
	}
	SendRawJson(w, data)
}
