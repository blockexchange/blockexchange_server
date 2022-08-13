package pages

import (
	"blockexchange/core"

	"github.com/gorilla/mux"
)

func (ctrl *Controller) SetupRoutes(r *mux.Router, cfg *core.Config) {
	r.HandleFunc("/", ctrl.Index)
	r.HandleFunc("/login", ctrl.Login)
}
