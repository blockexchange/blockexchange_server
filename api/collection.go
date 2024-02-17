package api

import (
	"blockexchange/types"
	"encoding/json"
	"fmt"
	"net/http"
)

func (api *Api) CreateCollection(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	c := &types.Collection{}
	err := json.NewDecoder(r.Body).Decode(c)
	if err != nil {
		SendError(w, 500, fmt.Sprintf("json error: %s", err))
		return
	}

	existing_collection, err := api.Repositories.CollectionRepo.GetCollectionByUserUIDAndName(c.UserUID, c.Name)
	if err != nil {
		SendError(w, 500, fmt.Sprintf("get existing collection error: %s", err))
		return
	}
	if existing_collection != nil {
		SendError(w, 405, fmt.Sprintf("collection with name '%s' already exists", c.Name))
		return
	}

	c.UserUID = claims.UserUID
	err = api.Repositories.CollectionRepo.CreateCollection(c)
	Send(w, c, err)
}
