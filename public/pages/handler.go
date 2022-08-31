package pages

import (
	"blockexchange/types"
	"errors"
	"net/http"
)

type RenderContext struct {
	baseUrl string
	ctrl    *Controller
	w       http.ResponseWriter
	r       *http.Request
}

func (rc *RenderContext) Render(file string, data any) error {
	return rc.ctrl.te.Execute("pages/about.html", rc.w, rc.r, rc.baseUrl, data)
}

type RenderFunc func(rc *RenderContext, r *http.Request) error

func (ctrl *Controller) Handler(baseUrl string, rf RenderFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rc := &RenderContext{
			ctrl:    ctrl,
			w:       w,
			r:       r,
			baseUrl: baseUrl,
		}

		err := rf(rc, r)
		if err != nil {
			ctrl.te.ExecuteError(w, r, "./", 500, err)
		}
	}
}

type SecureRenderFunc func(rc *RenderContext, r *http.Request, claims *types.Claims) error

func (ctrl *Controller) SecureHandler(baseUrl string, shf SecureRenderFunc, req_perms ...types.JWTPermission) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := ctrl.te.GetClaims(r)
		if err != nil {
			ctrl.te.ExecuteError(w, r, baseUrl, 500, err)
			return
		}
		if c == nil {
			ctrl.te.ExecuteError(w, r, baseUrl, 401, errors.New("unauthorized"))
			return
		}
		for _, req_perm := range req_perms {
			if !c.HasPermission(req_perm) {
				ctrl.te.ExecuteError(w, r, baseUrl, 403, errors.New("forbidden"))
				return
			}
		}

		rc := &RenderContext{
			ctrl:    ctrl,
			w:       w,
			r:       r,
			baseUrl: baseUrl,
		}

		err = shf(rc, r, c)
		if err != nil {
			ctrl.te.ExecuteError(w, r, baseUrl, 500, err)
			return
		}
	}
}
