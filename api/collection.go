package api

import (
	"blockexchange/core"
	"blockexchange/types"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func (api *Api) CreateOrUpdateCollection(w http.ResponseWriter, r *http.Request, ctx *SecureContext) {
	c := &types.Collection{}
	err := json.NewDecoder(r.Body).Decode(c)
	if err != nil {
		SendError(w, 500, fmt.Sprintf("json error: %s", err))
		return
	}

	if !core.ValidateName(c.Name) {
		SendError(w, 405, fmt.Sprintf("invalid name '%s'", c.Name))
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

	if c.UID == "" {
		// create
		c.UserUID = ctx.Claims.UserUID
		err = api.Repositories.CollectionRepo.CreateCollection(c)
		Send(w, c, err)
	} else {
		// update

		// fetch existing collection
		existing_collection, err = api.Repositories.CollectionRepo.GetCollectionByUID(c.UID)
		if err != nil {
			SendError(w, 500, fmt.Sprintf("get existing collection error (update): %s", err))
			return
		}
		if existing_collection == nil {
			SendError(w, 404, fmt.Sprintf("collection '%s' not found", c.UID))
			return
		}
		if existing_collection.UserUID != ctx.Claims.UserUID {
			SendError(w, 403, fmt.Sprintf("not allowed to modify collection '%s', owned by '%s'", c.UID, existing_collection.UserUID))
			return
		}

		// update allowed fields
		existing_collection.Name = c.Name
		existing_collection.Description = c.Description

		err = api.Repositories.CollectionRepo.UpdateCollection(existing_collection)
		Send(w, c, err)
	}
}

func (api *Api) DeleteCollection(w http.ResponseWriter, r *http.Request, ctx *SecureContext) {
	vars := mux.Vars(r)
	collection_uid := vars["collection_uid"]

	c, err := api.CollectionRepo.GetCollectionByUID(collection_uid)
	if err != nil {
		SendError(w, 500, fmt.Sprintf("get existing collection error: %s", err))
		return
	}
	if c == nil {
		SendError(w, 404, fmt.Sprintf("collection not found '%s'", collection_uid))
		return
	}
	if c.UserUID != ctx.Claims.UserUID {
		SendError(w, 403, fmt.Sprintf("not allowed to modify collection '%s', owned by '%s'", c.UID, c.UserUID))
		return
	}

	err = api.Repositories.CollectionRepo.DeleteCollection(collection_uid)
	Send(w, true, err)
}

func (api *Api) GetCollection(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	collection_uid := vars["collection_uid"]

	c, err := api.CollectionRepo.GetCollectionByUID(collection_uid)
	if err != nil {
		SendError(w, 500, fmt.Sprintf("get existing collection error: %s", err))
		return
	}
	if c == nil {
		SendError(w, 404, fmt.Sprintf("collection not found '%s'", collection_uid))
		return
	}

	Send(w, c, err)
}
