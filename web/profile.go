package web

import (
	"blockexchange/core"
	"blockexchange/db"
	"blockexchange/types"
	"blockexchange/web/components"
	"errors"
	"html/template"
	"net/http"
	"time"

	"github.com/gorilla/csrf"
)

type ProfileModel struct {
	UpdateError error
	User        *types.User
	AccessToken *components.AccessTokenModel
	CSRFField   template.HTML
}

func updateProfileData(m *ProfileModel, ur db.UserRepository, r *http.Request, c *types.Claims) error {
	var err error
	m.User, err = ur.GetUserById(c.UserID)
	if err != nil {
		return err
	}

	newUsername := r.Form.Get("username")
	newMail := r.Form.Get("mail")

	if newUsername != m.User.Name {
		if !core.ValidateName(newUsername) {
			m.UpdateError = errors.New("invalid username, allowed chars: [a-zA-Z0-9_.-]")
			return nil
		}

		newUser, err := ur.GetUserByName(newUsername)
		if err != nil {
			return err
		}
		if newUser != nil {
			m.UpdateError = errors.New("username already taken")
			return nil
		}

		m.User.Name = newUsername
	}

	m.User.Mail = &newMail
	return ur.UpdateUser(m.User)
}

func (ctx *Context) Profile(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	m := &ProfileModel{
		AccessToken: components.AccessToken(&ctx.Repos.AccessTokenRepo, r, c),
		CSRFField:   csrf.TemplateField(r),
	}

	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			ctx.RenderError(w, r, 500, err)
			return
		}
		if r.FormValue("action") == "update_profile" {
			m.User, err = ctx.Repos.UserRepo.GetUserById(c.UserID)
			if err != nil {
				ctx.RenderError(w, r, 500, err)
				return
			}

			err = updateProfileData(m, ctx.Repos.UserRepo, r, c)
			if err != nil {
				ctx.RenderError(w, r, 500, err)
				return
			}

			permissions := core.GetPermissions(m.User, true)
			dur := time.Duration(24 * 180 * time.Hour)
			token, err := core.CreateJWT(m.User, permissions, dur)
			if err != nil {
				ctx.RenderError(w, r, 500, err)
				return
			}

			ctx.SetClaims(w, token, dur)
		}
	}

	if m.User == nil {
		var err error
		m.User, err = ctx.Repos.UserRepo.GetUserById(c.UserID)
		if err != nil {
			ctx.RenderError(w, r, 500, err)
			return
		}
		if m.User == nil {
			ctx.RenderError(w, r, 404, errors.New("not found"))
			return
		}
	}

	ctx.ExecuteTemplate(w, r, "profile.html", m)
}
