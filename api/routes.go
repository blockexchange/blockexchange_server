package api

import (
	"blockexchange/core"

	"github.com/gorilla/mux"
)

func (api *Api) SetupRoutes(r *mux.Router, cfg *core.Config) {

	// api surface
	r.Handle("/api/info", InfoHandler{Config: cfg})
	r.HandleFunc("/api/token", api.RequestToken).Methods("POST")

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

}
