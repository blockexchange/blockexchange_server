package api

import (
	"blockexchange/types"
	"encoding/json"
	"net/http"
)

func (api *Api) Register(w http.ResponseWriter, r *http.Request) {
	rr := &types.RegisterRequest{}
	err := json.NewDecoder(r.Body).Decode(rr)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	_, resp, err := api.core.RegisterLocal(rr)
	Send(w, resp, err)
}
