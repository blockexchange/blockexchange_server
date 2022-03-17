package web

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (api *Api) GetCollectionSchemaByCollectionID(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	collection_id, err := strconv.Atoi(vars["collection_id"])
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	list, err := api.CollectionSchemaRepository.GetByCollectionID(int64(collection_id))
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}
	SendJson(w, list)
}

func (api *Api) CreateCollectionSchema(w http.ResponseWriter, r *http.Request, ctx *SecureContext) {
	vars := mux.Vars(r)
	collection_id, err := strconv.Atoi(vars["collection_id"])
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	schema_id, err := strconv.Atoi(vars["schema_id"])
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	schema, err := api.SchemaRepo.GetSchemaById(int64(schema_id))
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	if schema == nil {
		SendError(w, 404, "schema not found")
		return
	}

	if schema.UserID != ctx.Token.UserID {
		SendError(w, 403, "not the schema owner")
		return
	}

	collection, err := api.CollectionRepo.GetByID(int64(collection_id))
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	if collection == nil {
		SendError(w, 404, "collection not found")
		return
	}

	if collection.UserID != ctx.Token.UserID {
		SendError(w, 403, "not the collection owner")
		return
	}

	err = api.CollectionSchemaRepository.Create(int64(collection_id), int64(schema_id))

	Send(w, nil, err)
}

//TODO: Delete
