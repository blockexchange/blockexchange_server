package web

import (
	"fmt"
	"net/http"
	"time"

	"blockexchange/api"
	"blockexchange/controller"
	"blockexchange/core"
	"blockexchange/public"
	"blockexchange/public/pages"

	"github.com/dchest/captcha"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Serve(db_ *sqlx.DB, cfg *core.Config) error {

	r := mux.NewRouter()
	r.Use(prometheusMiddleware)
	r.Use(loggingMiddleware)

	// cache/store setup
	var cache core.Cache = core.NewNoOpCache()
	captchaExp := 10 * time.Minute
	var captchaStore captcha.Store = captcha.NewMemoryStore(50, captchaExp)

	if cfg.RedisHost != "" && cfg.RedisPort != "" {
		rdb := redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort),
			Password: "",
			DB:       0,
		})

		cache = core.NewRedisCache(rdb)
		captchaStore = core.NewRedisCaptchaStore(rdb, captchaExp)
	}
	captcha.SetCustomStore(captchaStore)

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
