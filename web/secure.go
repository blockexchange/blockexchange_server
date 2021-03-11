package web

import "net/http"

type SecureContext struct {
}

type Handler func(w http.ResponseWriter, r *http.Request)
type SecureHandler func(w http.ResponseWriter, r *http.Request, ctx *SecureContext)

func Secure(h SecureHandler) Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO
		h(w, r, nil)
	}
}
