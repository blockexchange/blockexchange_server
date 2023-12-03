package schema

import (
	"blockexchange/db"
	"blockexchange/tmpl"
	"blockexchange/types"

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
	r.HandleFunc("/schema/{username}/{schemaname}/delete", ctx.tu.Secure(ctx.SchemaDelete, tmpl.PermissionCheck(types.JWTPermissionManagement)))
	r.HandleFunc("/schema/{username}/{schemaname}/edit", ctx.tu.Secure(ctx.SchemaEdit, tmpl.PermissionCheck(types.JWTPermissionManagement)))
	r.HandleFunc("/schema/{username}/{schemaname}/edit-tags", ctx.tu.Secure(ctx.SchemaTagEdit, tmpl.PermissionCheck(types.JWTPermissionManagement)))
	r.HandleFunc("/schema/{username}", ctx.tu.OptionalSecure(ctx.UserSchema))
}
