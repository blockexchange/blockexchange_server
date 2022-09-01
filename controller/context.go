package controller

import (
	"blockexchange/core"
	"blockexchange/db"
	"blockexchange/types"
	"net/http"
	"time"
)

type RenderContext struct {
	baseUrl string
	ctrl    *Controller
	w       http.ResponseWriter
	r       *http.Request
	claims  *types.Claims
}

func (rc *RenderContext) Render(file string, data any) error {
	rd := &RenderData{
		BaseURL: rc.baseUrl,
		Data:    data,
	}

	c, err := rc.ctrl.GetClaims(rc.r)
	if err != nil {
		return err
	}

	rd.Claims = c

	if c != nil {
		rd.IsAdmin = c.HasPermission(types.JWTPermissionAdmin)
	}

	return rc.ctrl.te.Execute(file, rc.w, rc.r, 200, rd)
}

func (rc *RenderContext) Repositories() *db.Repositories {
	return rc.ctrl.Repositories
}

func (rc *RenderContext) Config() *core.Config {
	return rc.ctrl.cfg
}

func (rc *RenderContext) BaseURL() string {
	return rc.baseUrl
}

func (rc *RenderContext) ResponseWriter() http.ResponseWriter {
	return rc.w
}

func (rc *RenderContext) Request() *http.Request {
	return rc.r
}

func (rc *RenderContext) Claims() *types.Claims {
	return rc.claims
}

func (rc *RenderContext) SetToken(t string, dur time.Duration) {
	rc.ctrl.SetToken(rc.w, t, dur)
}

func (rc *RenderContext) RemoveToken() {
	rc.ctrl.RemoveToken(rc.w)
}
