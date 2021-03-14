package web

import (
	"blockexchange/db"
	"blockexchange/types"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type SchemaApi struct {
	SchemaRepo db.SchemaRepository
}

func DecodeSchema(rc io.ReadCloser) *types.Schema {
	m := make(map[string]interface{})
	json.NewDecoder(rc).Decode(&m)
	schema := types.Schema{}

	schema.Name = m["name"].(string)
	schema.Description = m["description"].(string)
	schema.MaxX = int(m["max_x"].(float64))

	return &schema
}

func (api SchemaApi) GetSchema(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	schema, err := api.SchemaRepo.GetSchemaById(int64(id))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(schema)
}

func (api SchemaApi) CreateSchema(w http.ResponseWriter, r *http.Request, ctx *SecureContext) {
	logrus.WithFields(logrus.Fields{
		"body": r.Body,
	}).Trace("POST /api/schema")
	if !ctx.CheckPermission(w, types.JWTPermissionUpload) {
		return
	}
	schema := DecodeSchema(r.Body)

	schema.UserID = ctx.Token.UserID

	err := api.SchemaRepo.CreateSchema(schema)
	if err != nil {
		SendError(w, err.Error())
		return
	}

	SendJson(w, schema)
}
