package pages

import (
	"blockexchange/controller"
	"blockexchange/core"
	"blockexchange/public/pages/schema"
	"blockexchange/types"

	"github.com/gorilla/mux"
)

func SetupRoutes(ctrl *controller.Controller, r *mux.Router, cfg *core.Config) {
	r.HandleFunc("/schema/{username}/{schemaname}/edit", ctrl.Handler("../../../", schema.SchemaEdit, types.JWTPermissionManagement))
	r.HandleFunc("/schema/{username}/{schemaname}/delete", ctrl.Handler("../../../", schema.SchemaDelete, types.JWTPermissionManagement))
	r.HandleFunc("/schema/{username}/{schemaname}/edit-tags", ctrl.Handler("../../../", schema.SchemaTagEdit, types.JWTPermissionManagement))

}
