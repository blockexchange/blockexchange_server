package pages

import (
	"blockexchange/core"

	"github.com/gorilla/mux"
)

func (ctrl *Controller) SetupRoutes(r *mux.Router, cfg *core.Config) {
	r.HandleFunc("/", ctrl.Index)
	r.HandleFunc("/login", ctrl.Login)
	r.HandleFunc("/schema/{username}/{schemaname}", ctrl.Schema)
	r.HandleFunc("/users", ctrl.Users)
	r.HandleFunc("/search", ctrl.Search)
	r.HandleFunc("/mod", ctrl.Mod)
	r.HandleFunc("/profile", ctrl.Profile)
	r.HandleFunc("/tags", ctrl.Profile)
	r.HandleFunc("/about", ctrl.About)
}