package web

import (
	"embed"
	"log"
	"net/http"
	"os"

	"blockexchange/public"

	"github.com/gorilla/mux"
)

func Serve() {

	// webdev flag
	useLocalfs := os.Getenv("WEBDEV") == "true"
	r := mux.NewRouter()

	// api surface
	r.HandleFunc("/api/info", InfoEndpoint)
	r.HandleFunc("/api/oauth_callback/github", OauthGithub)
	r.HandleFunc("/api/schema/{id}", GetSchema).Methods("GET")
	r.HandleFunc("/api/access_token", Secure(GetAccessTokens)).Methods("GET")
	r.HandleFunc("/api/access_token", Secure(PostAccessToken)).Methods("POST")
	r.HandleFunc("/api/access_token/{id}", Secure(DeleteAccessToken)).Methods("DELETE")

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
