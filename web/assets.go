package web

import (
	"embed"
	"net/http"

	"github.com/vearutop/statigz"
	"github.com/vearutop/statigz/brotli"
)

func HandleAssets(assets embed.FS, cache bool) http.Handler {
	if cache {
		return statigz.FileServer(assets, brotli.AddEncoding)
	} else {
		return http.FileServer(http.FS(assets))
	}
}
