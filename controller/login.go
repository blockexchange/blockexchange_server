package controller

import (
	"blockexchange/core"
	"blockexchange/types"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type LoginModel struct {
	Username   string
	Password   string
	LoginError bool
	Config     *core.Config
}

func (ctrl *Controller) Login(w http.ResponseWriter, r *http.Request) {
	lm := &LoginModel{
		Config: ctrl.cfg,
	}

	if r.Method == http.MethodPost {
		r.ParseForm()

		switch r.Form.Get("action") {
		case "login":
			lm.Username = r.Form.Get("username")
			lm.Password = r.Form.Get("password")

			user, err := ctrl.UserRepo.GetUserByName(lm.Username)
			if err != nil {
				ctrl.te.ExecuteError(w, r, "./", 500, err)
				return
			}

			if user == nil {
				lm.LoginError = true
				w.WriteHeader(401)
				break
			}

			err = bcrypt.CompareHashAndPassword([]byte(user.Hash), []byte(lm.Password))
			if err != nil {
				lm.LoginError = true
				w.WriteHeader(401)
				break
			}

			permissions := core.GetPermissions(user, true)
			if user.Role == types.UserRoleAdmin {
				// admin user
				permissions = append(permissions, types.JWTPermissionAdmin)
			}

			dur := time.Duration(7 * 24 * time.Hour)
			token, err := core.CreateJWT(user, permissions, dur)
			if err != nil {
				ctrl.te.ExecuteError(w, r, "./", 500, err)
				return
			}

			ctrl.te.SetToken(w, token, dur)
			http.Redirect(w, r, "login", http.StatusSeeOther)

		case "logout":
			ctrl.te.RemoveToken(w)
			http.Redirect(w, r, "login", http.StatusSeeOther)
			return
		}
	}

	ctrl.te.Execute("pages/login.html", w, r, "./", lm)
}
