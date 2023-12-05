package api

import (
	"blockexchange/types"
	"encoding/json"
	"net/http"
)

type SecureContext struct {
	Claims *types.Claims
}

func (ctx *SecureContext) CheckPermission(w http.ResponseWriter, permission types.JWTPermission) bool {
	for _, p := range ctx.Claims.Permissions {
		if p == permission {
			return true
		}
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
		claims, err := a.core.GetClaims(r)
		if err != nil {
			SendError(w, http.StatusForbidden, err.Error())
			return
		}
		if claims == nil {
			SendError(w, http.StatusUnauthorized, err.Error())
			return
		}

		ctx := SecureContext{
			Claims: claims,
		}

		h(w, r, &ctx)
	}
}
