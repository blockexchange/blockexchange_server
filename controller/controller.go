package controller

import (
	"blockexchange/core"
	"blockexchange/db"
	"blockexchange/public"
	"blockexchange/templateengine"
	"io/fs"
	"os"

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

	var f fs.FS = public.Files
	if cfg.WebDev {
		// dev mod: use local files instead of bundled
		f = os.DirFS("public")
	}

	return &Controller{
		Repositories: db.NewRepositories(db_),
		cfg:          cfg,
		te: templateengine.NewTemplateEngine(&templateengine.TemplateEngineOptions{
			Templates:   f,
			EnableCache: !cfg.WebDev,
			FuncMap:     funcs,
		}),
	}
}

func (ctrl *Controller) TemplateEngine() *templateengine.TemplateEngine {
	return ctrl.te
}
