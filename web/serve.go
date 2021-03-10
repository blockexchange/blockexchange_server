package web

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"
)

func Serve(content embed.FS) {

	// webdev flag
	useLocalfs := os.Getenv("WEBDEV") == "true"
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(getFileSystem(useLocalfs, content)))
	mux.HandleFunc("/api/info", InfoEndpoint)
	http.Handle("/", mux)

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
