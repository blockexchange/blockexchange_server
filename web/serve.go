package web

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"
)

func Serve(content embed.FS) {

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(getFileSystem(false, content)))
	http.Handle("/", mux)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

func getFileSystem(useOS bool, content embed.FS) http.FileSystem {
	if useOS {
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
