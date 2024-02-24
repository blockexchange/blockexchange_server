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
	"github.com/gorilla/mux"
	"github.com/minetest-go/colormapping"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/vearutop/statigz"
	"github.com/vearutop/statigz/brotli"
)

type Api struct {
	*db.Repositories
	cfg          *types.Config
	core         *core.Core
	Cache        *core.RedisCache
	ColorMapping *colormapping.ColorMapping
	Running      *atomic.Bool
}

func (a *Api) Stop() {
	a.Running.Store(false)
}

func NewApi(repos *db.Repositories, cfg *types.Config, rdb *redis.Client) (*Api, *mux.Router, error) {
	r := mux.NewRouter()
	r.Use(middleware.PrometheusMiddleware)
	r.Use(middleware.LoggingMiddleware)

	// cache/store setup
	captchaExp := 10 * time.Minute

	if cfg.RedisHost == "" || cfg.RedisPort == "" {
		return nil, nil, fmt.Errorf("redis not configured")
	}

	cache := core.NewRedisCache(rdb)
	captchaStore := core.NewRedisCaptchaStore(rdb, captchaExp)
	captcha.SetCustomStore(captchaStore)

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
		Cache:        cache,
		ColorMapping: cm,
		Running:      running,
	}

	// api surface
	a.SetupRoutes(r, cfg)
	// client app routes
	a.SetupSPARoutes(r, cfg)

	// static files
	if cfg.WebDev {
		logrus.WithFields(logrus.Fields{"dir": "public"}).Info("Using live mode")
		fs := http.FileServer(http.FS(os.DirFS("public")))
		r.PathPrefix("/").HandlerFunc(fs.ServeHTTP)

	} else {
		logrus.Info("Using embed mode")
		r.PathPrefix("/").Handler(statigz.FileServer(public.Webapp, brotli.AddEncoding))
	}

	return a, r, nil
}

func init() {
	// metrics
	http.Handle("/metrics", promhttp.Handler())
}
