package pages

import (
	"blockexchange/types"
	"errors"
	"net/http"
)

type SecureHandler func(http.ResponseWriter, *http.Request, *types.Claims)

func (ctrl *Controller) Secure(baseUrl string, h SecureHandler, req_perms ...types.JWTPermission) http.HandlerFunc {
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
		h(w, r, c)
	}
}
