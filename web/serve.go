package web

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"blockexchange/api"
	"blockexchange/core"
	"blockexchange/tmpl"

	"github.com/dchest/captcha"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func prettysize(num int) string {
	if num > (1000 * 1000) {
		return fmt.Sprintf("%d MB", num/(1000*1000))
	} else if num > 1000 {
		return fmt.Sprintf("%d kB", num/(1000))
	} else {
		return fmt.Sprintf("%d bytes", num)
	}
}

func formattime(ts int64) string {
	t := time.UnixMilli(ts)
	return t.Format(time.UnixDate)
}

func Serve(db_ *sqlx.DB, cfg *core.Config) error {

	r := mux.NewRouter()
	r.Use(prometheusMiddleware)
	r.Use(loggingMiddleware)
	r.Use(csrf.Protect([]byte(cfg.Key)))

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

	tu := &tmpl.TemplateUtil{
		Files: Files,
		AddFuncs: func(funcs template.FuncMap, r *http.Request) {
			funcs["BaseURL"] = func() string { return cfg.BaseURL }
			funcs["prettysize"] = prettysize
			funcs["formattime"] = formattime
		},
		JWTKey:       cfg.Key,
		CookieName:   cfg.CookieName,
		CookieDomain: cfg.CookieDomain,
		CookiePath:   cfg.CookiePath,
		CookieSecure: cfg.CookieSecure,
	}

	// templates, pages
	ctx := &Context{
		tu:     tu,
		Config: cfg,
		Repos:  a.Repositories,
	}
	ctx.Setup(r)

	// main entry
	http.Handle("/", r)

	// metrics
	http.Handle("/metrics", promhttp.Handler())

	return http.ListenAndServe(":8080", nil)
}
