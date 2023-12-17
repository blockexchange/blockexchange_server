package api

import (
	"blockexchange/api/middleware"
	"blockexchange/core"
	"blockexchange/db"
	"blockexchange/public"
	"blockexchange/types"
	"fmt"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	"github.com/dchest/captcha"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"github.com/minetest-go/colormapping"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"github.com/vearutop/statigz"
	"github.com/vearutop/statigz/brotli"

	"github.com/jmoiron/sqlx"
)

type Api struct {
	*db.Repositories
	cfg          *types.Config
	core         *core.Core
	Cache        core.Cache
	ColorMapping *colormapping.ColorMapping
	Running      *atomic.Bool
}

func (a *Api) Stop() {
	a.Running.Store(false)
}

func NewApi(db_ *sqlx.DB, cfg *types.Config) (*Api, error) {
	r := mux.NewRouter()
	r.Use(middleware.PrometheusMiddleware)
	r.Use(middleware.LoggingMiddleware)

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

	cm := colormapping.NewColorMapping()
	err := cm.LoadDefaults()
	if err != nil {
		return nil, err
	}

	running := &atomic.Bool{}
	running.Store(true)

	repos := db.NewRepositories(db_)

	a := &Api{
		Repositories: repos,
		cfg:          cfg,
		core:         core.New(cfg, repos),
		Cache:        cache,
		ColorMapping: cm,
		Running:      running,
	}

	// api surface
	a.SetupRoutes(r, cfg)
	// client app routes
	a.SetupWebRoutes(r, cfg)

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
