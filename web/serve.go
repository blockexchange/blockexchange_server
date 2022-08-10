package web

import (
	"net/http"
	"os"

	"blockexchange/core"
	"blockexchange/public"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Serve(db_ *sqlx.DB, cfg *core.Config) error {

	r := mux.NewRouter()

	// cache
	redis_host := os.Getenv("REDIS_HOST")
	redis_port := os.Getenv("REDIS_PORT")
	var cache core.Cache
	if redis_host != "" && redis_port != "" {
		cache = core.NewRedisCache(redis_host + ":" + redis_port)
	} else {
		cache = core.NewNoOpCache()
	}

	api, err := NewApi(db_, cache)
	if err != nil {
		return err
	}
	SetupRoutes(r, api, cfg)
	http.Handle("/", r)

	// metrics
	http.Handle("/metrics", promhttp.Handler())

	return http.ListenAndServe(":8080", nil)
}

func SetupRoutes(r *mux.Router, api *Api, cfg *core.Config) {

	// api surface
	r.Handle("/api/info", InfoHandler{Config: cfg})
	r.HandleFunc("/api/token", api.PostLogin).Methods("POST")

	r.HandleFunc("/api/export_we/{id}/{filename}", api.ExportWorldeditSchema).Methods("GET")
	r.HandleFunc("/api/export_bx/{id}/{filename}", api.ExportBXSchema).Methods("GET")

	r.HandleFunc("/api/schema/{id}", api.GetSchema).Methods("GET")
	r.HandleFunc("/api/schema", Secure(api.CreateSchema)).Methods("POST")
	r.HandleFunc("/api/schema/{id}", Secure(api.UpdateSchema)).Methods("PUT")
	r.HandleFunc("/api/schema/{id}/mods", api.GetSchemaMods).Methods("GET")
	r.HandleFunc("/api/schema/{id}/mods", Secure(api.CreateSchemaMods)).Methods("POST")
	r.HandleFunc("/api/schema/{id}/update", Secure(api.UpdateSchemaInfo)).Methods("POST")

	r.HandleFunc("/api/schema/{schema_id}/screenshot/{id}", api.GetSchemaScreenshotByID)
	r.HandleFunc("/api/schema/{schema_id}/screenshot", api.GetSchemaScreenshots)

	r.HandleFunc("/api/search/schema/byname/{user_name}/{schema_name}", api.SearchSchemaByNameAndUser)
	r.HandleFunc("/api/searchschema", api.SearchSchema).Methods("POST")

	r.HandleFunc("/api/schemapart", Secure(api.CreateSchemaPart)).Methods("POST")
	r.HandleFunc("/api/schemapart/{schema_id}/{x}/{y}/{z}", api.GetSchemaPart).Methods("GET")
	r.HandleFunc("/api/schemapart/{schema_id}/{x}/{y}/{z}", Secure(api.DeleteSchemaPart)).Methods("DELETE")
	r.HandleFunc("/api/schemapart/{schema_id}/{x}/{y}/{z}/delete", Secure(api.DeleteSchemaPart)).Methods("POST")
	r.HandleFunc("/api/schemapart_chunk/{schema_id}/{x}/{y}/{z}", api.GetSchemaPartChunk)
	r.HandleFunc("/api/schemapart_next/{schema_id}/{x}/{y}/{z}", api.GetNextSchemaPart)
	r.HandleFunc("/api/schemapart_next/by-mtime/{schema_id}/{mtime}", api.GetNextSchemaPartByMtime)
	r.HandleFunc("/api/schemapart_first/{schema_id}", api.GetFirstSchemaPart)

	r.HandleFunc("/", api.Index)

	r.PathPrefix("/assets/").HandlerFunc(HandleAssets(public.Files, webdev != "true"))

}
