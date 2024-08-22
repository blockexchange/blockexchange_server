package api

import (
	"blockexchange/schematic/bx"
	"blockexchange/schematic/worldedit"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/hashicorp/golang-lru/v2/expirable"
)

var ip_downloads_cache = expirable.NewLRU[string, bool](1000, nil, time.Hour*10)

func (api *Api) incrementDownloadStats(schema_uid string, r *http.Request) error {
	origin := r.Header.Get("X-Forwarded-For")
	origin_parts := strings.Split(origin, ",")
	if len(origin_parts) >= 1 {
		ip := strings.TrimSpace(origin_parts[0])
		cache_key := fmt.Sprintf("download_marker/%s/%s", schema_uid, ip)

		m, ok := ip_downloads_cache.Get(cache_key)
		if ok && m {
			//already incremented stat for this ip, skip
			return nil
		}

		// mark as incremented
		ip_downloads_cache.Add(cache_key, true)
	}

	err := api.SchemaRepo.IncrementDownloads(schema_uid)
	if err != nil {
		return err
	}
	return nil
}

var ip_views_cache = expirable.NewLRU[string, bool](1000, nil, time.Hour*10)

func (api *Api) incrementViewStats(schema_uid string, r *http.Request) error {
	origin := r.Header.Get("X-Forwarded-For")
	origin_parts := strings.Split(origin, ",")
	if len(origin_parts) >= 1 {
		ip := strings.TrimSpace(origin_parts[0])
		cache_key := fmt.Sprintf("view_marker/%s/%s", schema_uid, ip)

		m, ok := ip_views_cache.Get(cache_key)
		if ok && m {
			//already incremented stat for this ip, skip
			return nil
		}

		// mark as incremented
		ip_views_cache.Add(cache_key, true)
	}

	err := api.SchemaRepo.IncrementViews(schema_uid)
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

	err = api.incrementDownloadStats(schema_uid, r)
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

	err = api.incrementDownloadStats(schema_uid, r)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}
}
