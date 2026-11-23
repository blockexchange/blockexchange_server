package api

import (
	"blockexchange/types"
	"encoding/json"
	"net/http"
)

type SecureContext struct {
	Claims *types.Claims
	// raw jwt of the request
	Token string
}

func (ctx *SecureContext) HasPermission(permission types.JWTPermission) bool {
	for _, p := range ctx.Claims.Permissions {
		if p == permission {
			return true
		}
	}

	return false
}

func (ctx *SecureContext) CheckPermission(w http.ResponseWriter, permission types.JWTPermission) bool {
	if ctx.HasPermission(permission) {
		return true
	}

	// no permission found
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusForbidden)
	json.NewEncoder(w).Encode(types.ErrorResponse{Message: "permission required: " + string(permission)})
	return false
}

type Handler func(w http.ResponseWriter, r *http.Request)
type SecureHandler func(w http.ResponseWriter, r *http.Request, ctx *SecureContext)

func (a *Api) Secure(h SecureHandler) Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := a.core.GetToken(r)
		if err != nil {
			SendError(w, http.StatusForbidden, err.Error())
			return
		}
		if token == "" {
			SendError(w, http.StatusUnauthorized, "unauthorized")
			return
		}
		claims, err := a.core.ParseJWT(token)
		if err != nil {
			SendError(w, http.StatusForbidden, err.Error())
			return
		}

		ctx := SecureContext{
			Claims: claims,
			Token:  token,
		}

		h(w, r, &ctx)
	}
}
