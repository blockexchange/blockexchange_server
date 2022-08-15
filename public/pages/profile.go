package pages

import (
	"blockexchange/core"
	"blockexchange/types"
	"errors"
	"net/http"
)

type ProfileModel struct {
	UpdateError string
}

func (ctrl *Controller) updateProfileData(r *http.Request, c *types.Claims) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}

	user, err := ctrl.UserRepo.GetUserById(c.UserID)
	if err != nil {
		return err
	}

	newUsername := r.Form.Get("username")
	newMail := r.Form.Get("mail")

	if newUsername != user.Name {
		if !core.ValidateName(newUsername) {
			return errors.New("invalid username")
		}

		newUser, err := ctrl.UserRepo.GetUserByName(newUsername)
		if err != nil {
			return err
		}
		if newUser != nil {
			return errors.New("username already taken")
		}

		user.Name = newUsername
	}

	user.Mail = &newMail
	return ctrl.UserRepo.UpdateUser(user)
}

func (ctrl *Controller) Profile(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	baseUrl := "./"

	if r.Method == http.MethodPost {
		err := ctrl.updateProfileData(r, c)
		if err != nil {
			ctrl.te.ExecuteError(w, r, baseUrl, 500, err)
			return
		}
	}

	ctrl.te.Execute("pages/profile.html", w, r, baseUrl, nil)
}
