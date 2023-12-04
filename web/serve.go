package web

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"blockexchange/api"
	"blockexchange/core"
	"blockexchange/public"
	"blockexchange/types"

	"github.com/dchest/captcha"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/vearutop/statigz"
	"github.com/vearutop/statigz/brotli"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Serve(db_ *sqlx.DB, cfg *types.Config) (*api.Api, error) {

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
	a, err := api.NewApi(db_, cache, cfg)
	if err != nil {
		return nil, err
	}
	a.SetupRoutes(r, cfg)

	// index.html or /
	r.HandleFunc("/", public.RenderIndex)
	r.HandleFunc("/index.html", public.RenderIndex)
	r.HandleFunc("/login", public.RenderIndex)
	r.HandleFunc("/profile", public.RenderIndex) //TODO: better url handling

	// static files
	if cfg.WebDev {
		logrus.WithFields(logrus.Fields{"dir": "public"}).Info("Using live mode")
		fs := http.FileServer(http.FS(os.DirFS("public")))
		r.PathPrefix("/").HandlerFunc(fs.ServeHTTP)

	} else {
		logrus.Info("Using embed mode")
		r.PathPrefix("/").Handler(statigz.FileServer(public.Webapp, brotli.AddEncoding))
	}

	// main entry
	http.Handle("/", r)

	// metrics
	http.Handle("/metrics", promhttp.Handler())

	return a, nil
}
