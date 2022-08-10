package templateengine

import (
	"blockexchange/core"
	"blockexchange/types"
	"net/http"
	"time"
)

func (te *TemplateEngine) createCookie(value string, d time.Duration) *http.Cookie {
	return &http.Cookie{
		Name:     te.options.CookieName,
		Value:    value,
		Path:     te.options.CookiePath,
		Expires:  time.Now().Add(d),
		Domain:   te.options.CookieDomain,
		HttpOnly: true,
		Secure:   te.options.CookieSecure,
	}
}

func (te *TemplateEngine) SetToken(w http.ResponseWriter, token string, d time.Duration) {
	http.SetCookie(w, te.createCookie(token, d))
}

func (te *TemplateEngine) GetToken(r *http.Request) string {
	c, err := r.Cookie(te.options.CookieName)
	if err != nil {
		return ""
	}

	return c.Value
}

func (te *TemplateEngine) RemoveToken(w http.ResponseWriter) {
	http.SetCookie(w, te.createCookie("", time.Duration(0)))
}

func (te *TemplateEngine) GetClaims(r *http.Request) (*types.Claims, error) {
	t := te.GetToken(r)
	if t == "" {
		return nil, nil
	}

	return core.ParseJWT(t)
}
