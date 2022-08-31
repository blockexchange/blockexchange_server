package controller

import (
	"blockexchange/core"
	"blockexchange/db"
	"net/http"
	"time"
)

type RenderContext struct {
	baseUrl string
	ctrl    *Controller
	w       http.ResponseWriter
	r       *http.Request
}

func (rc *RenderContext) Render(file string, data any) error {
	return rc.ctrl.te.Execute(file, rc.w, rc.r, rc.baseUrl, data)
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

func (rc *RenderContext) SetToken(t string, dur time.Duration) {
	rc.ctrl.te.SetToken(rc.w, t, dur)
}

func (rc *RenderContext) RemoveToken() {
	rc.ctrl.te.RemoveToken(rc.w)
}
