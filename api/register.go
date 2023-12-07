package api

import (
	"blockexchange/types"
	"encoding/json"
	"net/http"
)

func (api *Api) CheckRegister(w http.ResponseWriter, r *http.Request) {
	rr := &types.RegisterRequest{}
	err := json.NewDecoder(r.Body).Decode(rr)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	resp, err := api.core.CheckRegister(rr, types.UserTypeLocal)
	Send(w, resp, err)
}

func (api *Api) Register(w http.ResponseWriter, r *http.Request) {
	rr := &types.RegisterRequest{}
	err := json.NewDecoder(r.Body).Decode(rr)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	_, err = api.core.Register(rr, types.UserTypeLocal)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}
}
