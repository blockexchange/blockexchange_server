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

	access_token_api := AccessTokenApi{Repo: &accessTokenRepo}
	token_api := TokenApi{
		AccessTokenRepo: &accessTokenRepo,
		UserRepo:        &userRepo,
	}
	schemaApi := SchemaApi{SchemaRepo: &schemaRepo}
	userApi := UserApi{UserRepo: &userRepo}
	oauthApi := OauthGithubApi{
		AccessTokenRepo: &accessTokenRepo,
		UserRepo:        &userRepo,
	}

	// api surface
	r.HandleFunc("/api/info", InfoEndpoint)
	r.HandleFunc("/api/token", token_api.PostLogin).Methods("POST")
	r.HandleFunc("/api/oauth_callback/github", oauthApi.OauthGithub)

	r.HandleFunc("/api/validate_username", userApi.PostValidateUsername).Methods("POST")
	r.HandleFunc("/api/user", userApi.GetUsers).Methods("GET")
	r.HandleFunc("/api/user/{id}", Secure(userApi.UpdateUser)).Methods("POST")

	r.HandleFunc("/api/schema/{id}", schemaApi.GetSchema).Methods("GET")

	r.HandleFunc("/api/access_token", Secure(access_token_api.GetAccessTokens)).Methods("GET")
	r.HandleFunc("/api/access_token", Secure(access_token_api.PostAccessToken)).Methods("POST")
	r.HandleFunc("/api/access_token/{id}", Secure(access_token_api.DeleteAccessToken)).Methods("DELETE")

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
