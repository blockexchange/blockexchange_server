package api

import (
	"blockexchange/schematic/bx"
	"blockexchange/schematic/worldedit"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

func (api *Api) incrementDownloadstats(schema_uid string, r *http.Request) error {
	origin := r.Header.Get("X-Forwarded-For")
	origin_parts := strings.Split(origin, ",")
	if len(origin_parts) >= 1 {
		ip := strings.TrimSpace(origin_parts[0])
		cache_key := fmt.Sprintf("download_marker/%s/%s", schema_uid, ip)

		m, err := api.Cache.Get(cache_key)
		if err != nil {
			return err
		}
		if m != nil {
			//already incremented stat for this ip, skip
			return nil
		}
		// mark as incremented
		err = api.Cache.Set(cache_key, []byte{0x00}, time.Hour*24)
		if err != nil {
			return err
		}
	}

	err := api.SchemaRepo.IncrementDownloads(schema_uid)
	if err != nil {
		return err
	}
	return nil
}

func (api *Api) ExportWorldeditSchema(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	schema_uid := vars["schema_uid"]

	w.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", vars["filename"]))
	e := worldedit.NewExporter(w)

	err := api.core.SchemapartCallback(schema_uid, e.Export)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	err = e.Close()
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	err = api.incrementDownloadstats(schema_uid, r)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}
}

func (api *Api) ExportBXSchema(w http.ResponseWriter, r *http.Request) {
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

	schemamods, err := api.SchemaModRepo.GetSchemaModsBySchemaUID(schema_uid)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	e := bx.NewExporter(w)
	err = e.ExportMetadata(schema, schemamods)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	err = api.core.SchemapartCallback(schema_uid, e.Export)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	err = e.Close()
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	err = api.incrementDownloadstats(schema_uid, r)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}
}
