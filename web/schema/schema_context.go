package schema

import (
	"blockexchange/db"
	"blockexchange/tmpl"

	"github.com/gorilla/mux"
)

type SchemaContext struct {
	tu      *tmpl.TemplateUtil
	repos   db.Repositories
	BaseURL string
}

func NewSchemaContext(tu *tmpl.TemplateUtil, repos db.Repositories, BaseURL string) *SchemaContext {
	return &SchemaContext{
		tu:      tu,
		repos:   repos,
		BaseURL: BaseURL,
	}
}

func (ctx *SchemaContext) Setup(r *mux.Router) {
	r.HandleFunc("/schema/{username}/{schemaname}", ctx.tu.OptionalSecure(ctx.Schema))

}
