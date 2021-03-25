package web

import (
	"blockexchange/types"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (api *Api) GetCollectionsByUserID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	user_id, err := strconv.Atoi(vars["user_id"])
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	list, err := api.CollectionRepo.GetByUserID(int64(user_id))
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	SendJson(w, list)
}

func (api *Api) CreateCollection(w http.ResponseWriter, r *http.Request, ctx *SecureContext) {
	collection := types.Collection{}
	err := json.NewDecoder(r.Body).Decode(&collection)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	if collection.UserID != ctx.Token.UserID {
		SendError(w, 403, "Userid does not match")
		return
	}

	err = api.CollectionRepo.Create(&collection)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	SendJson(w, collection)
}

func (api *Api) UpdateCollection(w http.ResponseWriter, r *http.Request, ctx *SecureContext) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	collection := types.Collection{}
	err = json.NewDecoder(r.Body).Decode(&collection)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	if id != int(collection.ID) {
		SendError(w, 405, "id does not match")
		return
	}

	if collection.UserID != ctx.Token.UserID {
		SendError(w, 403, "Userid does not match")
		return
	}

	err = api.CollectionRepo.Create(&collection)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	SendJson(w, collection)

}

func (api *Api) DeleteCollection(w http.ResponseWriter, r *http.Request, ctx *SecureContext) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	collection, err := api.CollectionRepo.GetByID(int64(id))
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	if collection.UserID != ctx.Token.UserID {
		SendError(w, 403, "Userid does not match")
		return
	}

	err = api.CollectionRepo.Delete(int64(id))
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}
