package web

import (
	"blockexchange/render"
	"blockexchange/types"
	"encoding/json"
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

	schema.Complete = true
	schema.Created = time.Now().Unix() * 1000
	err = api.SchemaRepo.UpdateSchema(schema)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	cm, err := render.GetColorMapping()
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	renderer := render.NewRenderer(api.SchemaPartRepo, cm)
	png, err := renderer.RenderSchema(schema)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	screenshots, err := api.SchemaScreenshotRepo.GetBySchemaID(schema.ID)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	var screenshot types.SchemaScreenshot

	if len(screenshots) >= 1 {
		// update existing
		screenshot = screenshots[0]

		err = api.SchemaScreenshotRepo.Update(&screenshot)
		if err != nil {
			SendError(w, 500, err.Error())
			return
		}
	} else {
		// create a new one
		screenshot = types.SchemaScreenshot{
			SchemaID: schema.ID,
			Type:     "image/png",
			Title:    "Isometric preview",
			Data:     png,
		}

		err = api.SchemaScreenshotRepo.Create(&screenshot)
		if err != nil {
			SendError(w, 500, err.Error())
			return
		}
	}

	err = api.SchemaRepo.CalculateStats(schema.ID)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}
