package web

import (
	"embed"
	"log"
	"net/http"
	"os"

	"blockexchange/public"
	"blockexchange/web/oauth"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

func Serve(db_ *sqlx.DB) {

	// webdev flag
	useLocalfs := os.Getenv("WEBDEV") == "true"
	r := mux.NewRouter()

	api := NewApi(db_)

	github_oauth := oauth.NewHandler(&oauth.GithubOauth{}, api.UserRepo, api.AccessTokenRepo)

	// api surface
	r.HandleFunc("/api/info", InfoEndpoint)
	r.HandleFunc("/api/token", api.PostLogin).Methods("POST")
	r.HandleFunc("/api/oauth_callback/github", github_oauth.Handle)

	r.HandleFunc("/api/validate_username", api.PostValidateUsername).Methods("POST")
	r.HandleFunc("/api/user", api.GetUsers).Methods("GET")
	r.HandleFunc("/api/user/{id}", Secure(api.UpdateUser)).Methods("POST")

	r.HandleFunc("/api/schema/{id}", api.GetSchema).Methods("GET")
	r.HandleFunc("/api/schema", Secure(api.CreateSchema)).Methods("POST")
	r.HandleFunc("/api/schema/{id}/mods", api.GetSchemaMods).Methods("GET")
	r.HandleFunc("/api/schema/{id}/mods", Secure(api.CreateSchemaMods)).Methods("POST")
	r.HandleFunc("/api/schema/{id}/complete", Secure(api.CompleteSchema)).Methods("POST")

	r.HandleFunc("/api/schema/{schema_id}/screenshot/{id}", api.GetSchemaScreenshotByID)

	r.HandleFunc("/api/search/schema/byname/{user_name}/{schema_name}", api.SearchSchemaByNameAndUser)
	r.HandleFunc("/api/searchschema", api.SearchSchema).Methods("POST")
	r.HandleFunc("/api/searchrecent/{count}", api.SearchRecentSchemas)

	r.HandleFunc("/api/schemapart", Secure(api.CreateSchemaPart)).Methods("POST")
	r.HandleFunc("/api/schemapart/{schema_id}/{x}/{y}/{z}", api.GetSchemaPart)
	r.HandleFunc("/api/schemapart_next/{schema_id}/{x}/{y}/{z}", api.GetNextSchemaPart)

	r.HandleFunc("/api/access_token", Secure(api.GetAccessTokens)).Methods("GET")
	r.HandleFunc("/api/access_token", Secure(api.PostAccessToken)).Methods("POST")
	r.HandleFunc("/api/access_token/{id}", Secure(api.DeleteAccessToken)).Methods("DELETE")

	// static files
	r.PathPrefix("/").Handler(http.FileServer(getFileSystem(useLocalfs, public.Webapp)))
	http.Handle("/", r)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

func getFileSystem(useLocalfs bool, content embed.FS) http.FileSystem {
	if useLocalfs {
		log.Print("using live mode")
		return http.FS(os.DirFS("public"))
	}

	log.Print("using embed mode")
	return http.FS(content)
}
