package controller

import (
	"blockexchange/types"
	"errors"
	"net/http"
)

type RenderFunc func(rc *RenderContext) error

func (ctrl *Controller) Handler(baseUrl string, shf RenderFunc, req_perms ...types.JWTPermission) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := ctrl.GetClaims(r)
		if err != nil {
			ctrl.te.Execute("pages/error.html", w, r, 500, &RenderData{BaseURL: baseUrl, Data: err})
			return
		}
		for _, req_perm := range req_perms {
			if !c.HasPermission(req_perm) {
				ctrl.te.Execute("pages/error.html", w, r, 403, &RenderData{BaseURL: baseUrl, Data: errors.New("forbidden")})
				return
			}
		}

		rc := &RenderContext{
			ctrl:               ctrl,
			w:                  w,
			r:                  r,
			baseUrl:            baseUrl,
			claims:             c,
			AdditionalMetaTags: make(map[string]string),
		}

		err = shf(rc)
		if err != nil {
			err = ctrl.te.Execute("pages/error.html", w, r, 500, &RenderData{BaseURL: baseUrl, Data: err})
			if err != nil {
				w.WriteHeader(500)
				w.Write([]byte(err.Error()))
			}
			return
		}
	}
}
