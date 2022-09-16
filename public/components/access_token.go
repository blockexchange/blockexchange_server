package components

import (
	"blockexchange/controller"
	"blockexchange/core"
	"blockexchange/types"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type AccessTokenModel struct {
	Err      error
	Message  string
	Tokens   []*types.AccessToken
	Username string
}

func handleAccessTokenAdd(rc *controller.RenderContext, m *AccessTokenModel) error {
	name := rc.Request().FormValue("name")
	expire_days_str := rc.Request().FormValue("expire_days")

	if name == "" {
		return errors.New("token name is empty")
	}
	expire_days, err := strconv.ParseInt(expire_days_str, 10, 64)
	if err != nil {
		return errors.New("Invalid expiration value: " + err.Error())
	}

	t := &types.AccessToken{
		UserID:  rc.Claims().UserID,
		Name:    name,
		Expires: (time.Now().Unix() + (expire_days * 3600 * 24)) * 1000,
		Created: time.Now().Unix() * 1000,
		Token:   core.CreateToken(6),
	}
	err = rc.Repositories().AccessTokenRepo.CreateAccessToken(t)
	if err != nil {
		return err
	}

	m.Message = fmt.Sprintf("New access-token '%s' created", t.Name)
	return nil
}

func handleAccessTokenRemove(rc *controller.RenderContext, m *AccessTokenModel) error {
	token_id_str := rc.Request().FormValue("token_id")
	token_name := rc.Request().FormValue("token_name")
	token_id, err := strconv.ParseInt(token_id_str, 10, 64)
	if err != nil {
		return err
	}

	err = rc.Repositories().AccessTokenRepo.RemoveAccessToken(token_id, rc.Claims().UserID)
	if err != nil {
		return err
	}

	m.Message = fmt.Sprintf("Access-token '%s' removed", token_name)
	return nil
}

func AccessToken(rc *controller.RenderContext) *AccessTokenModel {
	m := &AccessTokenModel{
		Username: rc.Claims().Username,
	}
	r := rc.Request()
	if r.Method == http.MethodPost {
		m.Err = r.ParseForm()
		if m.Err != nil {
			return m
		}
		if r.FormValue("action") == "add_token" {
			m.Err = handleAccessTokenAdd(rc, m)
			if m.Err != nil {
				return m
			}
		}
		if r.FormValue("action") == "remove_token" {
			m.Err = handleAccessTokenRemove(rc, m)
			if m.Err != nil {
				return m
			}
		}
	}

	m.Tokens, m.Err = rc.Repositories().AccessTokenRepo.GetAccessTokensByUserID(rc.Claims().UserID)
	return m
}
