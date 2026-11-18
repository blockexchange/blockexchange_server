package api

import (
	"blockexchange/api/middleware"
	"blockexchange/core"
	"blockexchange/db"
	"blockexchange/public"
	"blockexchange/types"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/gorilla/mux"
	"github.com/luanti-golang/colormapping"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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

	// static files
	if cfg.WebDev {
		logrus.WithFields(logrus.Fields{"dir": "public/dist"}).Info("Using live mode")
		fs := http.FileServer(http.FS(os.DirFS("public/dist")))
		r.PathPrefix("/").HandlerFunc(fs.ServeHTTP)

	} else {
		logrus.Info("Using embed mode")
		r.PathPrefix("/").Handler(statigz.FileServer(public.Webapp, brotli.AddEncoding, statigz.FSPrefix(public.Prefix)))
	}

	return a, r, nil
}

func init() {
	// metrics
	http.Handle("/metrics", promhttp.Handler())
}
