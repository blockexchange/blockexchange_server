package api

import (
	"blockexchange/schematic/bx"
	"blockexchange/schematic/worldedit"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/hashicorp/golang-lru/v2/expirable"
)

var ip_downloads_cache = expirable.NewLRU[string, bool](1000, nil, time.Hour*10)
var ip_views_cache = expirable.NewLRU[string, bool](1000, nil, time.Hour*10)

// returns the client ip, preferring the first X-Forwarded-For entry
// and falling back to the connection's remote address
func clientIP(r *http.Request) string {
	forwarded := strings.TrimSpace(strings.Split(r.Header.Get("X-Forwarded-For"), ",")[0])
	if forwarded != "" {
		return forwarded
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

// marks the (schema, client ip) pair in the cache and returns true
// if it was not marked yet, i.e. the stat should be counted
func markFirstHit(cache *expirable.LRU[string, bool], schema_uid string, r *http.Request) bool {
	cache_key := fmt.Sprintf("%s/%s", schema_uid, clientIP(r))
	if marked, ok := cache.Get(cache_key); ok && marked {
		return false
	}
	cache.Add(cache_key, true)
	return true
}

func (api *Api) incrementDownloadStats(schema_uid string, r *http.Request) error {
	if !markFirstHit(ip_downloads_cache, schema_uid, r) {
		return nil
	}
	return api.SchemaRepo.IncrementDownloads(schema_uid)
}

func (api *Api) incrementViewStats(schema_uid string, r *http.Request) error {
	if !markFirstHit(ip_views_cache, schema_uid, r) {
		return nil
	}
	return api.SchemaRepo.IncrementViews(schema_uid)
}

func (api *Api) ExportWorldeditSchema(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	schema_uid := vars["schema_uid"]

	w.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", vars["filename"]))
	e := worldedit.NewExporter(w)

	err := api.core.SchemapartCallback(schema_uid, e.Export)
	if SendIfError(w, err) {
		return
	}

	err = e.Close()
	if SendIfError(w, err) {
		return
	}

	err = api.incrementDownloadStats(schema_uid, r)
	if SendIfError(w, err) {
		return
	}
}

func (api *Api) ExportBXSchema(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	schema_uid := vars["schema_uid"]

	schema, err := api.SchemaRepo.GetSchemaByUID(schema_uid)
	if SendIfError(w, err) {
		return
	}
	if schema == nil {
		SendError(w, 404, "not found")
		return
	}

	schemamods, err := api.SchemaModRepo.GetSchemaModsBySchemaUID(schema_uid)
	if SendIfError(w, err) {
		return
	}

	e := bx.NewExporter(w)
	err = e.ExportMetadata(schema, schemamods)
	if SendIfError(w, err) {
		return
	}

	err = api.core.SchemapartCallback(schema_uid, e.Export)
	if SendIfError(w, err) {
		return
	}

	err = e.Close()
	if SendIfError(w, err) {
		return
	}

	err = api.incrementDownloadStats(schema_uid, r)
	if SendIfError(w, err) {
		return
	}
}
