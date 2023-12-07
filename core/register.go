package core

import (
	"blockexchange/types"
	"fmt"
	"time"

	"github.com/dchest/captcha"
	"golang.org/x/crypto/bcrypt"
)

func (c *Core) Register(rr *types.RegisterRequest, ut types.UserType) (*types.User, error) {
	resp, err := c.CheckRegister(rr, ut)
	if err != nil {
		return nil, err
	}
	if !resp.Success {
		return nil, fmt.Errorf("register error")
	}

	hash := []byte("invalid")

	if ut == types.UserTypeGithub {
		hash, err = bcrypt.GenerateFromPassword([]byte(rr.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
	}

	user := &types.User{
		Created: time.Now().Unix() * 1000,
		Name:    rr.Name,
		Type:    ut,
		Role:    types.UserRoleDefault,
		Hash:    string(hash),
		Mail:    &rr.Mail,
	}
	err = c.repos.UserRepo.CreateUser(user)
	if err != nil {
		return nil, err
	}

	err = c.repos.AccessTokenRepo.CreateAccessToken(&types.AccessToken{
		Name:    "default",
		Created: time.Now().Unix() * 1000,
		Expires: (time.Now().Unix() + (3600 * 24 * 7 * 4)) * 1000,
		Token:   CreateToken(6),
		UserID:  *user.ID,
	})
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (c *Core) CheckRegister(rr *types.RegisterRequest, ut types.UserType) (*types.CheckRegisterResponse, error) {
	resp := &types.CheckRegisterResponse{}
	if !ValidateName(rr.Name) || rr.Name == "" {
		resp.ErrInvalidUsername = true
		return resp, nil
	}

	existing_user, err := c.repos.UserRepo.GetUserByName(rr.Name)
	if err != nil {
		return nil, err
	}
	if existing_user != nil {
		resp.ErrUsernameTaken = true
		return resp, nil
	}

	if ut == types.UserTypeLocal {
		// additional checks fr local user
		if len(rr.Password) < 6 {
			resp.ErrPasswordTooShort = true
			return resp, nil
		}

		if !captcha.VerifyString(rr.CaptchaID, rr.CaptchaAnswer) {
			resp.ErrCaptcha = true
			return resp, nil
		}
	}

	resp.Success = true
	return resp, nil
}
