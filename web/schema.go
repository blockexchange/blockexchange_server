package web

import (
	"blockexchange/core"
	"blockexchange/types"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func (api Api) GetSchema(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}
	schema, err := api.SchemaRepo.GetSchemaById(int64(id))
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	if r.URL.Query().Get("download") == "true" {
		schema.Downloads++
		err = api.SchemaRepo.UpdateSchema(schema)
		if err != nil {
			SendError(w, 500, err.Error())
			return
		}
	}

	SendJson(w, schema)
}

func (api Api) DeleteSchema(w http.ResponseWriter, r *http.Request, ctx *SecureContext) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	err = api.SchemaRepo.DeleteSchema(int64(id), ctx.Token.UserID)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (api Api) CreateSchema(w http.ResponseWriter, r *http.Request, ctx *SecureContext) {
	logrus.WithFields(logrus.Fields{
		"body": r.Body,
	}).Trace("POST /api/schema")
	if !ctx.CheckPermission(w, types.JWTPermissionUpload) {
		return
	}
	schema := types.Schema{}
	err := json.NewDecoder(r.Body).Decode(&schema)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	if !core.ValidateName(schema.Name) {
		SendError(w, 400, "Invalid schema name")
		return
	}

	MAX_SIZE := 500
	if schema.SizeX > MAX_SIZE || schema.SizeY > MAX_SIZE || schema.SizeZ > MAX_SIZE {
		SendError(w, 400, fmt.Sprintf("Max side-length of %d nodes exceeded", MAX_SIZE))
		return
	}

	// remove incomplete schema with same name if it exists
	err = api.SchemaRepo.DeleteIncompleteSchema(ctx.Token.UserID, schema.Name)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	schema.UserID = ctx.Token.UserID
	schema.Created = time.Now().Unix() * 1000

	err = api.SchemaRepo.CreateSchema(&schema)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	SendJson(w, schema)
}

func (api Api) UpdateSchema(w http.ResponseWriter, r *http.Request, ctx *SecureContext) {
	logrus.WithFields(logrus.Fields{
		"body": r.Body,
	}).Trace("PUT /api/schema")

	if !ctx.CheckPermission(w, types.JWTPermissionManagement) {
		return
	}

	schema := types.Schema{}
	err := json.NewDecoder(r.Body).Decode(&schema)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	schema.UserID = ctx.Token.UserID
	schema.Created = time.Now().Unix() * 1000

	err = api.SchemaRepo.UpdateSchema(&schema)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	SendJson(w, schema)
}

func (api Api) UpdateSchemaInfo(w http.ResponseWriter, r *http.Request, ctx *SecureContext) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	schema, err := api.SchemaRepo.GetSchemaById(int64(id))
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	if schema.UserID != ctx.Token.UserID {
		SendError(w, 403, "you are not the owner of the schema")
		return
	}

	// schema is complete now, mark as initial
	initial := !schema.Complete

	schema.Complete = true
	schema.Created = time.Now().Unix() * 1000
	err = api.SchemaRepo.UpdateSchema(schema)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	// update screenshot
	screenshot, err := core.UpdatePreview(schema, api.Repositories)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	// let the database calculate the size/count stats
	err = api.SchemaRepo.CalculateStats(schema.ID)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	// retrieve updated schema data from the db (size, count)
	schema, err = api.SchemaRepo.GetSchemaById(schema.ID)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	if initial {
		user, err := api.UserRepo.GetUserById(schema.UserID)
		if err != nil {
			SendError(w, 500, err.Error())
			return
		}

		// initial schema upload, send it to the feed async
		go core.UpdateSchemaFeed(schema, user, screenshot)
	}

	w.WriteHeader(http.StatusOK)
}
