package pages

import (
	"blockexchange/controller"
	"blockexchange/core"
	"blockexchange/db"
	"blockexchange/types"
	"errors"
	"net/http"
)

type ProfileModel struct {
	UpdateError string
}

func updateProfileData(ur db.UserRepository, r *http.Request, c *types.Claims) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}

	user, err := ur.GetUserById(c.UserID)
	if err != nil {
		return err
	}

	newUsername := r.Form.Get("username")
	newMail := r.Form.Get("mail")

	if newUsername != user.Name {
		if !core.ValidateName(newUsername) {
			return errors.New("invalid username")
		}

		newUser, err := ur.GetUserByName(newUsername)
		if err != nil {
			return err
		}
		if newUser != nil {
			return errors.New("username already taken")
		}

		user.Name = newUsername
	}

	user.Mail = &newMail
	return ur.UpdateUser(user)
}

func Profile(rc *controller.RenderContext, r *http.Request, c *types.Claims) error {
	if r.Method == http.MethodPost {
		err := updateProfileData(rc.Repositories().UserRepo, r, c)
		if err != nil {
			return err
		}
	}

	return rc.Render("pages/profile.html", nil)
}
