package api

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
	id, err := strconv.Atoi(vars["schema_id"])
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
		err = api.incrementDownloadstats(int64(id), r)
		if err != nil {
			SendError(w, 500, err.Error())
			return
		}
	}

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
	err = api.SchemaRepo.DeleteIncompleteSchema(ctx.Claims.UserID, schema.Name)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	schema.UserID = ctx.Claims.UserID
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

	updated_schema := types.Schema{}
	err := json.NewDecoder(r.Body).Decode(&updated_schema)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	// fetch saved schema
	schema, err := api.SchemaRepo.GetSchemaById(updated_schema.ID)
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
	if !is_admin && schema.UserID != ctx.Claims.UserID {
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
	existing_schema, err := api.SchemaRepo.GetSchemaByUserIDAndName(schema.UserID, updated_schema.Name)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}
	if existing_schema != nil && existing_schema.ID != schema.ID {
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
	err = api.SchemaRepo.UpdateSchema(schema)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	SendJson(w, schema)
}

func (api Api) UpdateSchemaInfo(w http.ResponseWriter, r *http.Request, ctx *SecureContext) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["schema_id"])
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	schema, err := api.SchemaRepo.GetSchemaById(int64(id))
	if err != nil {
		SendError(w, 500, fmt.Sprintf("GetSchemaById: %s", err))
		return
	}
	if schema == nil {
		SendError(w, 404, "not found")
		return
	}

	if !ctx.HasPermission(types.JWTPermissionAdmin) && schema.UserID != ctx.Claims.UserID {
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
	err = api.SchemaRepo.CalculateStats(schema.ID)
	if err != nil {
		SendError(w, 500, fmt.Sprintf("CalculateStats: %s", err))
		return
	}

	// retrieve updated schema data from the db (size, count)
	schema, err = api.SchemaRepo.GetSchemaById(schema.ID)
	if err != nil {
		SendError(w, 500, fmt.Sprintf("GetBySchemaID: %s", err))
		return
	}

	// process notifications
	if notify_feed {
		screenshots, err := api.SchemaScreenshotRepo.GetBySchemaID(schema.ID)
		if err != nil {
			SendError(w, 500, fmt.Sprintf("GetBySchemaID: %s", err))
			return
		}

		user, err := api.UserRepo.GetUserById(schema.UserID)
		if err != nil {
			SendError(w, 500, fmt.Sprintf("GetUserById: %s", err))
			return
		}

		if len(screenshots) > 0 {
			go core.UpdateSchemaFeed(schema, user, screenshots[0])
		}
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
	id, err := strconv.ParseInt(vars["schema_id"], 10, 64)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	schema, err := api.SchemaRepo.GetSchemaById(id)
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
	if !is_admin && schema.UserID != ctx.Claims.UserID {
		// not an admin and not the owner
		SendError(w, 403, "unauthorized")
		return
	}

	err = api.SchemaRepo.DeleteSchema(schema.ID)
	Send(w, true, err)
}
