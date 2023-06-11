package components

import (
	"blockexchange/core"
	"blockexchange/db"
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

func handleAccessTokenAdd(repo *db.AccessTokenRepository, r *http.Request, c *types.Claims, m *AccessTokenModel) error {
	name := r.FormValue("name")
	expire_days_str := r.FormValue("expire_days")

	if name == "" {
		return errors.New("token name is empty")
	}
	expire_days, err := strconv.ParseInt(expire_days_str, 10, 64)
	if err != nil {
		return errors.New("Invalid expiration value: " + err.Error())
	}

	t := &types.AccessToken{
		UserID:  c.UserID,
		Name:    name,
		Expires: (time.Now().Unix() + (expire_days * 3600 * 24)) * 1000,
		Created: time.Now().Unix() * 1000,
		Token:   core.CreateToken(6),
	}
	err = repo.CreateAccessToken(t)
	if err != nil {
		return err
	}

	m.Message = fmt.Sprintf("New access-token '%s' created", t.Name)
	return nil
}

func handleAccessTokenRemove(repo *db.AccessTokenRepository, r *http.Request, c *types.Claims, m *AccessTokenModel) error {
	token_id_str := r.FormValue("token_id")
	token_name := r.FormValue("token_name")
	token_id, err := strconv.ParseInt(token_id_str, 10, 64)
	if err != nil {
		return err
	}

	err = repo.RemoveAccessToken(token_id, c.UserID)
	if err != nil {
		return err
	}

	m.Message = fmt.Sprintf("Access-token '%s' removed", token_name)
	return nil
}

func AccessToken(repo *db.AccessTokenRepository, r *http.Request, c *types.Claims) *AccessTokenModel {
	m := &AccessTokenModel{
		Username: c.Username,
	}
	if r.Method == http.MethodPost {
		m.Err = r.ParseForm()
		if m.Err != nil {
			return m
		}
		if r.FormValue("action") == "add_token" {
			m.Err = handleAccessTokenAdd(repo, r, c, m)
			if m.Err != nil {
				return m
			}
		}
		if r.FormValue("action") == "remove_token" {
			m.Err = handleAccessTokenRemove(repo, r, c, m)
			if m.Err != nil {
				return m
			}
		}
	}

	m.Tokens, m.Err = repo.GetAccessTokensByUserID(c.UserID)
	return m
}
