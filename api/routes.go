package api

import (
	"blockexchange/core"
	"net/http"

	"github.com/dchest/captcha"
	"github.com/gorilla/mux"
)

func (api *Api) SetupRoutes(r *mux.Router, cfg *core.Config) {

	// api surface
	r.Handle("/api/info", InfoHandler{Config: cfg})
	r.HandleFunc("/api/token", api.RequestToken).Methods(http.MethodPost)
	r.PathPrefix("/api/captcha/").Handler(captcha.Server(350, 250))

	r.HandleFunc("/api/export_we/{id}/{filename}", api.ExportWorldeditSchema).Methods(http.MethodGet)
	r.HandleFunc("/api/export_bx/{id}/{filename}", api.ExportBXSchema).Methods(http.MethodGet)

	r.HandleFunc("/api/schema/{id}", api.GetSchema).Methods(http.MethodGet)
	r.HandleFunc("/api/schema", Secure(api.CreateSchema)).Methods(http.MethodPost)
	r.HandleFunc("/api/schema/{id}", Secure(api.UpdateSchema)).Methods(http.MethodPut)
	r.HandleFunc("/api/schema/{id}/mods", api.GetSchemaMods).Methods(http.MethodGet)
	r.HandleFunc("/api/schema/{id}/mods", Secure(api.CreateSchemaMods)).Methods(http.MethodPost)
	r.HandleFunc("/api/schema/{id}/update", Secure(api.UpdateSchemaInfo)).Methods(http.MethodPost)

	r.HandleFunc("/api/schema/{schema_id}/screenshot/update", Secure(api.UpdateSchemaPreview)).Methods(http.MethodPost)
	r.HandleFunc("/api/schema/{schema_id}/screenshot", api.GetFirstSchemaScreenshot)

	r.HandleFunc("/api/search/schema/byname/{user_name}/{schema_name}", api.SearchSchemaByNameAndUser)
	r.HandleFunc("/api/search/schema", api.SearchSchema).Methods(http.MethodPost)

	r.HandleFunc("/api/schemapart", Secure(api.CreateSchemaPart)).Methods(http.MethodPost)
	r.HandleFunc("/api/schemapart/{schema_id}/{x}/{y}/{z}", api.GetSchemaPart).Methods(http.MethodGet)
	r.HandleFunc("/api/schemapart/{schema_id}/{x}/{y}/{z}", Secure(api.DeleteSchemaPart)).Methods("DELETE")
	r.HandleFunc("/api/schemapart/{schema_id}/{x}/{y}/{z}/delete", Secure(api.DeleteSchemaPart)).Methods(http.MethodPost)
	r.HandleFunc("/api/schemapart_chunk/{schema_id}/{x}/{y}/{z}", api.GetSchemaPartChunk)
	r.HandleFunc("/api/schemapart_next/{schema_id}/{x}/{y}/{z}", api.GetNextSchemaPart)
	r.HandleFunc("/api/schemapart_next/by-mtime/{schema_id}/{mtime}", api.GetNextSchemaPartByMtime)
	r.HandleFunc("/api/schemapart_count/by-mtime/{schema_id}/{mtime}", api.CountNextSchemaPartByMtime)
	r.HandleFunc("/api/schemapart_first/{schema_id}", api.GetFirstSchemaPart)

}
