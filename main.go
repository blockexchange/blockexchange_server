package main

import (
	"blockexchange/web"
	"embed"
)

//go:embed public/*
var content embed.FS

func main() {
	println("Starting")
	web.Serve(content)
}
