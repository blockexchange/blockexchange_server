package controller

import (
	"blockexchange/templateengine"
	"net/http"
)

type LoginModel struct {
	Username   string
	Password   string
	LoginError bool
}

func (ctrl *Controller) Login(w http.ResponseWriter, r *http.Request) {
	lm := &LoginModel{}

	if r.Method == http.MethodPost {
		r.ParseForm()

		switch r.Form.Get("action") {
		case "login":
			lm.Username = r.Form.Get("username")
			lm.Password = r.Form.Get("password")

			if lm.Password == "enter" {
				// ok
				claims := &templateengine.Claims{
					Username: r.Form.Get("username"),
				}
				ctrl.te.SetClaims(w, claims)
				http.Redirect(w, r, "login", http.StatusSeeOther)

			} else {
				// wrong pw
				lm.LoginError = true
				w.WriteHeader(401)
			}
		case "logout":
			ctrl.te.RemoveClaims(w)
			http.Redirect(w, r, "login", http.StatusSeeOther)
			return
		}
	}

	ctrl.te.Execute("pages/login.html", w, r, "./", lm)
}
