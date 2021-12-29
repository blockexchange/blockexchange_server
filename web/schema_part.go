package web

import (
	"blockexchange/types"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var partsUploaded = promauto.NewCounter(prometheus.CounterOpts{
	Name: "bx_schemaparts_uploaded",
	Help: "The total number of uploaded schemaparts",
})

var partsDownloaded = promauto.NewCounter(prometheus.CounterOpts{
	Name: "bx_schemaparts_downloaded",
	Help: "The total number of downloaded schemaparts",
})

func (api *Api) CreateSchemaPart(w http.ResponseWriter, r *http.Request, ctx *SecureContext) {
	schemapart := types.SchemaPart{}
	err := json.NewDecoder(r.Body).Decode(&schemapart)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	schema, err := api.SchemaRepo.GetSchemaById(schemapart.SchemaID)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	if schema == nil {
		SendError(w, 500, "no schema found")
		return
	}

	if schema.UserID != ctx.Token.UserID {
		SendError(w, 403, "you are not the owner of the schema")
		return
	}

	mtime := time.Now().Unix() * 1000

	// update schema part
	schemapart.Mtime = mtime
	err = api.SchemaPartRepo.CreateOrUpdateSchemaPart(&schemapart)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	// update schema mtime
	schema.Mtime = mtime
	err = api.SchemaRepo.UpdateSchema(schema)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	// increment stats
	partsUploaded.Inc()

	SendJson(w, schemapart)
}

func extractSchemaPartVars(r *http.Request) (int64, int, int, int, error) {
	vars := mux.Vars(r)
	schema_id, err := strconv.Atoi(vars["schema_id"])
	if err != nil {
		return 0, 0, 0, 0, err
	}
	x, err := strconv.Atoi(vars["x"])
	if err != nil {
		return 0, 0, 0, 0, err
	}

	y, err := strconv.Atoi(vars["y"])
	if err != nil {
		return 0, 0, 0, 0, err
	}

	z, err := strconv.Atoi(vars["z"])
	if err != nil {
		return 0, 0, 0, 0, err
	}

	return int64(schema_id), x, y, z, nil
}

func (api *Api) GetSchemaPart(w http.ResponseWriter, r *http.Request) {
	schema_id, x, y, z, err := extractSchemaPartVars(r)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	partsDownloaded.Inc()

	schemapart, err := api.SchemaPartRepo.GetBySchemaIDAndOffset(int64(schema_id), x, y, z)
	if err == nil && schemapart == nil {
		w.WriteHeader(http.StatusNoContent)
	} else {
		Send(w, schemapart, err)
	}
}

func (api *Api) GetSchemaPartChunk(w http.ResponseWriter, r *http.Request) {
	schema_id, x, y, z, err := extractSchemaPartVars(r)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	to_pos_offset := 16 * 4
	schemaparts, err := api.SchemaPartRepo.GetBySchemaIDAndRange(int64(schema_id), x, y, z, x+to_pos_offset, y+to_pos_offset, z+to_pos_offset)
	if err == nil && schemaparts == nil {
		w.WriteHeader(http.StatusNoContent)
	} else {
		partsDownloaded.Add(float64(len(schemaparts)))
		Send(w, schemaparts, err)
	}
}

func (api *Api) GetNextSchemaPart(w http.ResponseWriter, r *http.Request) {
	schema_id, x, y, z, err := extractSchemaPartVars(r)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	partsDownloaded.Inc()

	schemapart, err := api.SchemaPartRepo.GetNextBySchemaIDAndOffset(int64(schema_id), x, y, z)
	if err == nil && schemapart == nil {
		w.WriteHeader(http.StatusNoContent)
	} else {
		Send(w, schemapart, err)
	}
}

func (api *Api) GetNextSchemaPartByMtime(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	schema_id, err := strconv.Atoi(vars["schema_id"])
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	mtime, err := strconv.Atoi(vars["mtime"])
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	partsDownloaded.Inc()

	schemapart, err := api.SchemaPartRepo.GetNextBySchemaIDAndMtime(int64(schema_id), int64(mtime))
	if err == nil && schemapart == nil {
		w.WriteHeader(http.StatusNoContent)
	} else {
		Send(w, schemapart, err)
	}
}

func (api *Api) GetFirstSchemaPart(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	schema_id, err := strconv.Atoi(vars["schema_id"])
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	partsDownloaded.Inc()

	schemapart, err := api.SchemaPartRepo.GetFirstBySchemaID(int64(schema_id))
	if err == nil && schemapart == nil {
		w.WriteHeader(http.StatusNoContent)
	} else {
		Send(w, schemapart, err)
	}
}
