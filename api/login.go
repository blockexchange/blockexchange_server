package api

import (
	"blockexchange/types"
	"encoding/json"
	"net/http"
)

func (api Api) Login(w http.ResponseWriter, r *http.Request) {
	login := types.Login{}
	err := json.NewDecoder(r.Body).Decode(&login)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	// TODO
}
