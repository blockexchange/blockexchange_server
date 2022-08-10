package controller

import (
	"blockexchange/core"
	"blockexchange/db"
	"blockexchange/public"
	"blockexchange/templateengine"
	"os"

	"github.com/jmoiron/sqlx"
)

type Controller struct {
	*db.Repositories
	cfg *core.Config
	te  *templateengine.TemplateEngine
}

func NewController(db_ *sqlx.DB, cfg *core.Config) *Controller {
	return &Controller{
		Repositories: db.NewRepositories(db_),
		cfg:          cfg,
		te: templateengine.NewTemplateEngine(&templateengine.TemplateEngineOptions{
			Templates:    public.Files,
			TemplateDir:  "public",
			EnableCache:  !cfg.WebDev,
			JWTKey:       os.Getenv("BLOCKEXCHANGE_KEY"),
			CookieName:   "blockexchange",
			CookiePath:   os.Getenv("BLOCKEXCHANGE_COOKIE_PATH"),
			CookieDomain: os.Getenv("BLOCKEXCHANGE_COOKIE_DOMAIN"),
			CookieSecure: os.Getenv("BLOCKEXCHANGE_COOKIE_SECURE") == "true",
		}),
	}
}
