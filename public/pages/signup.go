package pages

import (
	"blockexchange/controller"
	"blockexchange/core"
	"blockexchange/types"
	"errors"
	"net/http"
	"time"

	"github.com/dchest/captcha"
	"golang.org/x/crypto/bcrypt"
)

type SignupModel struct {
	CaptchaID string
	Username  string
	Mail      string
	Password  string
	Password2 string
	Err       error
}

func Signup(rc *controller.RenderContext) error {
	m := &SignupModel{
		CaptchaID: captcha.New(),
	}

	r := rc.Request()
	if r.Method == http.MethodPost {
		m.Err = handleSignup(rc, m)
	}

	return rc.Render("pages/signup.html", m)
}

func handleSignup(rc *controller.RenderContext, m *SignupModel) error {
	r := rc.Request()
	r.ParseForm()
	m.Username = r.FormValue("username")
	m.Password = r.FormValue("password")
	m.Password2 = r.FormValue("password2")
	m.Mail = r.FormValue("mail")
	m.CaptchaID = r.FormValue("captcha_id")
	entered_captcha := r.FormValue("captcha")

	if !core.ValidateName(m.Username) || m.Username == "" {
		return errors.New("invalid username, allowed characters: a-zA-Z0-9_.-")
	}

	existing_user, err := rc.Repositories().UserRepo.GetUserByName(m.Username)
	if err != nil {
		return err
	}

	if existing_user != nil {
		return errors.New("username already taken")
	}

	if len(m.Password) < 6 {
		return errors.New("password too short, should at least be 6 characters")
	}

	if m.Password != m.Password2 {
		return errors.New("passwords do not match")
	}

	if !captcha.VerifyString(m.CaptchaID, entered_captcha) {
		m.CaptchaID = captcha.New()
		return errors.New("captcha invalid")
	}

	//TODO: deduplicate (in handler.go)

	hash, err := bcrypt.GenerateFromPassword([]byte(m.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &types.User{
		Created: time.Now().Unix() * 1000,
		Name:    m.Username,
		Type:    types.UserTypeLocal,
		Hash:    string(hash),
		Mail:    &m.Mail,
	}
	err = rc.Repositories().UserRepo.CreateUser(user)
	if err != nil {
		return err
	}

	err = rc.Repositories().AccessTokenRepo.CreateAccessToken(&types.AccessToken{
		Name:    "default",
		Created: time.Now().Unix() * 1000,
		Expires: (time.Now().Unix() + (3600 * 24 * 7 * 4)) * 1000,
		Token:   core.CreateToken(6),
		UserID:  *user.ID,
	})
	if err != nil {
		return err
	}

	dur := time.Duration(24 * 180 * time.Hour)
	permissions := core.GetPermissions(user, true)
	token, err := core.CreateJWT(user, permissions, dur)
	if err != nil {
		return err
	}

	rc.SetToken(token, dur)
	rc.Redirect(rc.BaseURL() + "/profile")

	return nil
}
