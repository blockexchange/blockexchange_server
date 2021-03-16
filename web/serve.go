package web

import (
	"embed"
	"log"
	"net/http"
	"os"

	"blockexchange/db"
	"blockexchange/public"

	"github.com/gorilla/mux"
)

func Serve() {

	// webdev flag
	useLocalfs := os.Getenv("WEBDEV") == "true"
	r := mux.NewRouter()

	accessTokenRepo := db.DBAccessTokenRepository{DB: db.DB}
	userRepo := db.DBUserRepository{DB: db.DB}
	schemaRepo := db.DBSchemaRepository{DB: db.DB}

	api := Api{
		AccessTokenRepo: accessTokenRepo,
		UserRepo:        userRepo,
		SchemaRepo:      schemaRepo,
	}

	// api surface
	r.HandleFunc("/api/info", InfoEndpoint)
	r.HandleFunc("/api/token", api.PostLogin).Methods("POST")
	r.HandleFunc("/api/oauth_callback/github", api.OauthGithub)

	r.HandleFunc("/api/validate_username", api.PostValidateUsername).Methods("POST")
	r.HandleFunc("/api/user", api.GetUsers).Methods("GET")
	r.HandleFunc("/api/user/{id}", Secure(api.UpdateUser)).Methods("POST")

	r.HandleFunc("/api/schema/{id}", api.GetSchema).Methods("GET")
	r.HandleFunc("/api/schema", Secure(api.CreateSchema)).Methods("POST")

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
