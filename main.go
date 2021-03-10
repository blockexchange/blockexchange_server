package main

import (
	"blockexchange/db"
	"blockexchange/web"
	"embed"
)

//go:embed public/*
var content embed.FS

func main() {
	println("Starting")
	db.Init()
	db.Migrate()
	web.Serve(content)
}
