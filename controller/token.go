package controller

import (
	"blockexchange/core"
	"blockexchange/types"
	"net/http"
	"time"
)

func (c *Controller) createCookie(value string, d time.Duration) *http.Cookie {
	return &http.Cookie{
		Name:     c.cfg.CookieName,
		Value:    value,
		Path:     c.cfg.CookiePath,
		Expires:  time.Now().Add(d),
		Domain:   c.cfg.CookieDomain,
		HttpOnly: true,
		Secure:   c.cfg.CookieSecure,
	}
}

func (c *Controller) SetToken(w http.ResponseWriter, token string, d time.Duration) {
	http.SetCookie(w, c.createCookie(token, d))
}

func (c *Controller) GetToken(r *http.Request) string {
	co, err := r.Cookie(c.cfg.CookieName)
	if err != nil {
		return ""
	}

	return co.Value
}

func (c *Controller) RemoveToken(w http.ResponseWriter) {
	http.SetCookie(w, c.createCookie("", time.Duration(0)))
}

func (c *Controller) GetClaims(r *http.Request) (*types.Claims, error) {
	t := c.GetToken(r)
	if t == "" {
		return nil, nil
	}

	return core.ParseJWT(t)
}
