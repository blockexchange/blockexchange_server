package web

import (
	"blockexchange/colormapping"
	"blockexchange/core"
	"blockexchange/db"
	"blockexchange/public"
	"blockexchange/templateengine"
	"os"

	"github.com/jmoiron/sqlx"
)

var webdev = os.Getenv("WEBDEV")

type Api struct {
	*db.Repositories
	te           *templateengine.TemplateEngine
	Cache        core.Cache
	ColorMapping *colormapping.ColorMapping
}

func NewApi(db_ *sqlx.DB, cache core.Cache) (*Api, error) {
	cm := colormapping.NewColorMapping()
	err := cm.LoadDefaults()
	if err != nil {
		return nil, err
	}

	return &Api{
		Repositories: db.NewRepositories(db_),
		Cache:        cache,
		ColorMapping: cm,
		te: templateengine.NewTemplateEngine(&templateengine.TemplateEngineOptions{
			Templates:    public.Files,
			TemplateDir:  "public",
			EnableCache:  webdev != "true",
			JWTKey:       os.Getenv("BLOCKEXCHANGE_KEY"),
			CookieName:   "blockexchange",
			CookiePath:   os.Getenv("BLOCKEXCHANGE_COOKIE_PATH"),
			CookieDomain: os.Getenv("BLOCKEXCHANGE_COOKIE_DOMAIN"),
			CookieSecure: os.Getenv("BLOCKEXCHANGE_COOKIE_SECURE") == "true",
		}),
	}, nil
}
