package api

import (
	"blockexchange/api/middleware"
	"blockexchange/core"
	"blockexchange/db"
	"blockexchange/public"
	"blockexchange/types"
	"blockexchange/wasm"
	"net/http"
	"os"
	"strings"
	"sync/atomic"

	"github.com/gorilla/mux"
	"github.com/luanti-golang/colormapping"
	"github.com/sirupsen/logrus"
	"github.com/vearutop/statigz"
	"github.com/vearutop/statigz/brotli"
)

type Api struct {
	*db.Repositories
	cfg          *types.Config
	core         *core.Core
	ColorMapping *colormapping.ColorMapping
	Running      *atomic.Bool
}

func (a *Api) Stop() {
	a.Running.Store(false)
}

func NewApi(repos *db.Repositories, cfg *types.Config) (*Api, *mux.Router, error) {
	r := mux.NewRouter()
	r.Use(middleware.PrometheusMiddleware)
	r.Use(middleware.LoggingMiddleware)
	r.Use(middleware.WasmHeadersMiddleware)

	cm := colormapping.NewColorMapping()
	err := cm.LoadDefaults()
	if err != nil {
		return nil, nil, err
	}

	running := &atomic.Bool{}
	running.Store(true)

	a := &Api{
		Repositories: repos,
		cfg:          cfg,
		core:         core.New(cfg, repos),
		ColorMapping: cm,
		Running:      running,
	}

	// api surface
	a.SetupRoutes(r, cfg)
	// redirects from the old history-mode frontend urls
	a.SetupLegacyRedirects(r)

	// luanti wasm build for the "Edit Online" client, served from the same origin (required for the pthread workers)
	wasm_runtime := statigz.FileServer(wasm.Runtime, statigz.FSPrefix(wasm.Prefix))
	r.PathPrefix(middleware.WasmRuntimePrefix).Handler(http.StripPrefix(middleware.WasmRuntimePrefix, noDirListing(wasm_runtime)))

	// static files
	var webapp http.Handler
	if cfg.WebDev {
		logrus.WithFields(logrus.Fields{"dir": "public/dist"}).Info("Using live mode")
		webapp = http.FileServer(http.FS(os.DirFS("public/dist")))
	} else {
		logrus.Info("Using embed mode")
		webapp = statigz.FileServer(public.Webapp, brotli.AddEncoding, statigz.FSPrefix(public.Prefix))
	}

	// the webapp again as cross-origin isolated document (headers by the wasm middleware) for the "Edit Online" page
	r.Path(middleware.WasmPrefix).Handler(http.StripPrefix(strings.TrimSuffix(middleware.WasmPrefix, "/"), webapp))
	r.PathPrefix("/").Handler(webapp)

	return a, r, nil
}

// refuses directory listings of a file server
func noDirListing(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "" || strings.HasSuffix(r.URL.Path, "/") {
			http.NotFound(w, r)
			return
		}
		h.ServeHTTP(w, r)
	})
}
