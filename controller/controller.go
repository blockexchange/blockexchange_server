package controller

import (
	"blockexchange/core"
	"blockexchange/db"
	"blockexchange/public"
	"blockexchange/templateengine"

	"github.com/jmoiron/sqlx"
)

type Controller struct {
	*db.Repositories
	cfg *core.Config
	te  *templateengine.TemplateEngine
}

func NewController(db_ *sqlx.DB, cfg *core.Config) *Controller {
	funcs := make(map[string]any)
	funcs["prettysize"] = prettysize
	funcs["formattime"] = formattime

	return &Controller{
		Repositories: db.NewRepositories(db_),
		cfg:          cfg,
		te: templateengine.NewTemplateEngine(&templateengine.TemplateEngineOptions{
			Templates:   public.Files,
			TemplateDir: "public",
			EnableCache: !cfg.WebDev,
			FuncMap:     funcs,
		}),
	}
}

func (ctrl *Controller) TemplateEngine() *templateengine.TemplateEngine {
	return ctrl.te
}
