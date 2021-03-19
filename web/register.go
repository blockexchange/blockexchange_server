package web

import (
	"blockexchange/types"
	"encoding/json"
	"net/http"
)

func (api *Api) Register(w http.ResponseWriter, r *http.Request) {
	info := types.RegisterInfo{}
	err := json.NewDecoder(r.Body).Decode(&info)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

}
