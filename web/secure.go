package web

import (
	"blockexchange/types"
	"encoding/json"
	"net/http"
)

type SecureContext struct {
}

type Handler func(w http.ResponseWriter, r *http.Request)
type SecureHandler func(w http.ResponseWriter, r *http.Request, ctx *SecureContext)

func Secure(h SecureHandler) Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		authorization := r.Header.Get("Authorization")
		if authorization == "" {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(types.ErrorResponse{Message: "no jwt found"})
			return
		}

		// TODO
		h(w, r, nil)
	}
}
