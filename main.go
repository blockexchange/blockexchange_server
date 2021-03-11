package main

import (
	"blockexchange/db"
	"blockexchange/web"
)

func main() {
	println("Starting")
	db.Init()
	db.Migrate()
	web.Serve()
}
