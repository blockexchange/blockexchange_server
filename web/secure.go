package web

import (
	"blockexchange/core"
	"blockexchange/types"
	"encoding/json"
	"net/http"
)

type SecureContext struct {
	Token *types.TokenInfo
}

func (ctx *SecureContext) CheckPermission(w http.ResponseWriter, permission types.JWTPermission) bool {
	for _, p := range ctx.Token.Permissions {
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

func Secure(h SecureHandler) Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		authorization := r.Header.Get("Authorization")
		if authorization == "" {
			SendError(w, http.StatusUnauthorized, "no jwt found")
			return
		}

		token, err := core.ParseJWT(authorization)
		if err != nil {
			SendError(w, http.StatusForbidden, err.Error())
			return
		}

		ctx := SecureContext{
			Token: token,
		}

		h(w, r, &ctx)
	}
}
