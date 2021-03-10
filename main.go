package main

import (
	"blockexchange/db"
	"blockexchange/web"
	"embed"
)

//go:embed public/js/* public/pics/* public/index.html
//go:embed public/node_modules/bootstrap/dist/css/bootstrap.min.css
//go:embed public/node_modules/vue/dist/vue.min.js
//go:embed public/node_modules/vue-router/dist/vue-router.min.js
//go:embed public/node_modules/vue-i18n/dist/vue-i18n.min.js
//go:embed public/node_modules/@fortawesome/fontawesome-free/css/all.min.css
//go:embed public/node_modules/@fortawesome/fontawesome-free/webfonts/*
var webapp embed.FS

func main() {
	println("Starting")
	db.Init()
	db.Migrate()
	web.Serve(webapp)
}
