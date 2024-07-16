package api

import (
	"blockexchange/core"
	"blockexchange/types"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func (api Api) GetSchema(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	schema_uid := vars["schema_uid"]

	schema, err := api.SchemaRepo.GetSchemaByUID(schema_uid)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	if r.URL.Query().Get("download") == "true" {
		err = api.incrementDownloadStats(schema_uid, r)
		if err != nil {
			SendError(w, 500, err.Error())
			return
		}
	}

	// increment view stats
	api.incrementViewStats(schema.UID, r)

	SendJson(w, schema)
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
	err = api.SchemaRepo.DeleteIncompleteSchema(ctx.Claims.UserUID, schema.Name)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	schema.UserUID = ctx.Claims.UserUID
	schema.Created = time.Now().UnixMilli()

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

	updated_schema := types.Schema{}
	err := json.NewDecoder(r.Body).Decode(&updated_schema)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	// fetch saved schema
	schema, err := api.SchemaRepo.GetSchemaByUID(updated_schema.UID)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}
	if schema == nil {
		SendError(w, 404, "not found")
		return
	}

	// check permissions
	is_admin := ctx.HasPermission(types.JWTPermissionAdmin)
	if !is_admin && schema.UserUID != ctx.Claims.UserUID {
		// not an admin and not the owner
		SendError(w, 403, "unauthorized")
		return
	}

	// check name
	if !core.ValidateName(updated_schema.Name) {
		SendErrorResponse(w, http.StatusBadRequest, &types.SchemaUpdateError{
			NameInvalid: true,
		})
		return
	}

	// check if the name already exists
	existing_schema, err := api.SchemaRepo.GetSchemaByUserUIDAndName(schema.UserUID, updated_schema.Name)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}
	if existing_schema != nil && existing_schema.UID != schema.UID {
		// another schema with the same name already exists
		SendErrorResponse(w, http.StatusBadRequest, &types.SchemaUpdateError{
			NameTaken: true,
		})
		return
	}

	// apply modifiable fields
	schema.Name = updated_schema.Name
	schema.License = updated_schema.License
	schema.Description = updated_schema.Description
	schema.ShortDescription = updated_schema.ShortDescription
	schema.CDBCollection = updated_schema.CDBCollection
	schema.CollectionUID = updated_schema.CollectionUID

	err = api.SchemaRepo.UpdateSchema(schema)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	SendJson(w, schema)
}

func (api Api) UpdateSchemaInfo(w http.ResponseWriter, r *http.Request, ctx *SecureContext) {
	vars := mux.Vars(r)
	schema_uid := vars["schema_uid"]

	schema, err := api.SchemaRepo.GetSchemaByUID(schema_uid)
	if err != nil {
		SendError(w, 500, fmt.Sprintf("GetSchemaById: %s", err))
		return
	}
	if schema == nil {
		SendError(w, 404, "not found")
		return
	}

	if !ctx.HasPermission(types.JWTPermissionAdmin) && schema.UserUID != ctx.Claims.UserUID {
		SendError(w, 403, "you are not the owner of the schema")
		return
	}

	notify_feed := false
	if !schema.Complete {
		// initial upload, complete schema
		schema.Complete = true
		schema.Created = time.Now().Unix() * 1000
		// set notify
		notify_feed = true

		err = api.SchemaRepo.UpdateSchema(schema)
		if err != nil {
			SendError(w, 500, fmt.Sprintf("UpdateSchema: %s", err))
			return
		}

		// update screenshot
		_, err := api.core.UpdatePreview(schema)
		if err != nil {
			SendError(w, 500, fmt.Sprintf("UpdatePreview: %s", err))
			return
		}

	}

	// let the database calculate the size/count stats
	err = api.SchemaRepo.CalculateStats(schema.UID)
	if err != nil {
		SendError(w, 500, fmt.Sprintf("CalculateStats: %s", err))
		return
	}

	// retrieve updated schema data from the db (size, count)
	schema, err = api.SchemaRepo.GetSchemaByUID(schema.UID)
	if err != nil {
		SendError(w, 500, fmt.Sprintf("GetBySchemaID: %s", err))
		return
	}

	// process notifications
	if notify_feed {
		user, err := api.UserRepo.GetUserByUID(schema.UserUID)
		if err != nil {
			SendError(w, 500, fmt.Sprintf("GetUserById: %s", err))
			return
		}

		go core.UpdateSchemaFeed(schema, user)
	}

	Send(w, schema, nil)
}

func (api Api) DeleteSchema(w http.ResponseWriter, r *http.Request, ctx *SecureContext) {
	logrus.WithFields(logrus.Fields{
		"body": r.Body,
	}).Trace("DELETE /api/schema")

	if !ctx.CheckPermission(w, types.JWTPermissionManagement) {
		return
	}

	// fetch schema
	vars := mux.Vars(r)
	schema_uid := vars["schema_uid"]

	schema, err := api.SchemaRepo.GetSchemaByUID(schema_uid)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}
	if schema == nil {
		SendError(w, 404, "not found")
		return
	}

	// check permissions
	is_admin := ctx.HasPermission(types.JWTPermissionAdmin)
	if !is_admin && schema.UserUID != ctx.Claims.UserUID {
		// not an admin and not the owner
		SendError(w, 403, "unauthorized")
		return
	}

	err = api.SchemaRepo.DeleteSchema(schema.UID)
	Send(w, true, err)
}
