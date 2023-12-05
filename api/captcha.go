package api

import (
	"net/http"

	"github.com/dchest/captcha"
)

func (api *Api) CreateCaptcha(w http.ResponseWriter, r *http.Request) {
	Send(w, captcha.New(), nil)
}
