package web

import (
	"net/http"
	"os"

	"blockexchange/api"
	"blockexchange/controller"
	"blockexchange/core"
	"blockexchange/public"
	"blockexchange/public/pages"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Serve(db_ *sqlx.DB, cfg *core.Config) error {

	r := mux.NewRouter()
	r.Use(prometheusMiddleware)
	r.Use(loggingMiddleware)

	// cache
	redis_host := os.Getenv("REDIS_HOST")
	redis_port := os.Getenv("REDIS_PORT")
	var cache core.Cache
	if redis_host != "" && redis_port != "" {
		cache = core.NewRedisCache(redis_host + ":" + redis_port)
	} else {
		cache = core.NewNoOpCache()
	}

	// api setup and routing
	a, err := api.NewApi(db_, cache)
	if err != nil {
		return err
	}
	a.SetupRoutes(r, cfg)

	// controller setup and routing
	ctrl := controller.NewController(db_, cfg)
	pages.SetupRoutes(ctrl, r, cfg)

	// assets
	r.PathPrefix("/assets/").Handler(HandleAssets(public.Files, !cfg.WebDev))

	// main entry
	http.Handle("/", r)

	// metrics
	http.Handle("/metrics", promhttp.Handler())

	return http.ListenAndServe(":8080", nil)
}
