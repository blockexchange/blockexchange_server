package web

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func Serve(content embed.FS) {

	// webdev flag
	useLocalfs := os.Getenv("WEBDEV") == "true"
	r := mux.NewRouter()

	// api surface
	r.HandleFunc("/api/info", InfoEndpoint)
	r.HandleFunc("/api/schema/{id}", GetSchema)
	r.HandleFunc("/api/oauth_callback/github", OauthGithub)

	// static files
	r.PathPrefix("/").Handler(http.FileServer(getFileSystem(useLocalfs, content)))
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
	fsys, err := fs.Sub(content, "public")
	if err != nil {
		panic(err)
	}

	return http.FS(fsys)
}
