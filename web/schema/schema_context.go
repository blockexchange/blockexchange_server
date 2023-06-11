package schema

import (
	"blockexchange/db"
	"blockexchange/web"
)

type SchemaContext struct {
	tu      *web.TemplateUtil
	repos   db.Repositories
	BaseURL string
}
