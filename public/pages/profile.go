package pages

import (
	"blockexchange/controller"
	"blockexchange/core"
	"blockexchange/db"
	"blockexchange/public/components"
	"blockexchange/types"
	"errors"
	"net/http"
	"time"
)

type ProfileModel struct {
	UpdateError error
	User        *types.User
	AccessToken *components.AccessTokenModel
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

func Profile(rc *controller.RenderContext) error {
	m := &ProfileModel{
		AccessToken: components.AccessToken(rc),
	}

	r := rc.Request()
	if r.Method == http.MethodPost {
		err := rc.Request().ParseForm()
		if err != nil {
			return err
		}
		m.User, err = rc.Repositories().UserRepo.GetUserById(rc.Claims().UserID)
		if err != nil {
			return err
		}

		err = updateProfileData(m, rc.Repositories().UserRepo, r, rc.Claims())
		if err != nil {
			return err
		}

		permissions := core.GetPermissions(m.User, true)
		dur := time.Duration(24 * 180 * time.Hour)
		token, err := core.CreateJWT(m.User, permissions, dur)
		if err != nil {
			return err
		}
		rc.SetToken(token, dur)
	}

	if m.User == nil {
		var err error
		m.User, err = rc.Repositories().UserRepo.GetUserById(rc.Claims().UserID)
		if err != nil {
			return err
		}
		if m.User == nil {
			return errors.New("not found")
		}
	}

	return rc.Render("pages/profile.html", m)
}
