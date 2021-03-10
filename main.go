package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"
)

func main() {
	println("Starting")
	http.Handle("/", http.FileServer(getFileSystem(false)))

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

//go:embed public/*
var content embed.FS

func getFileSystem(useOS bool) http.FileSystem {
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
