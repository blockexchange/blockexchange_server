package templateengine

import (
	"errors"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var err_unauthorized = errors.New("unauthorized")

func (te *TemplateEngine) createCookie(value string) *http.Cookie {
	return &http.Cookie{
		Name:     te.options.CookieName,
		Value:    value,
		Path:     te.options.CookiePath,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		Domain:   te.options.CookieDomain,
		HttpOnly: true,
		Secure:   te.options.CookieSecure,
	}
}

func (te *TemplateEngine) SetToken(w http.ResponseWriter, token string) {
	http.SetCookie(w, te.createCookie(token))
}

func (te *TemplateEngine) GetToken(r *http.Request) string {
	c, err := r.Cookie(te.options.CookieName)
	if err != nil {
		return ""
	}

	return c.Value
}

func (te *TemplateEngine) RemoveClaims(w http.ResponseWriter) {
	http.SetCookie(w, te.createCookie(""))
}

func (te *TemplateEngine) SetClaims(w http.ResponseWriter, claims *Claims) error {
	claims.RegisteredClaims = &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 12)),
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err := t.SignedString([]byte(te.options.JWTKey))
	if err != nil {
		return err
	}

	te.SetToken(w, token)
	return nil
}

func (te *TemplateEngine) GetClaims(r *http.Request) (*Claims, error) {
	t := te.GetToken(r)
	if t == "" {
		return nil, nil
	}

	token, err := jwt.ParseWithClaims(t, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(te.options.JWTKey), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, err_unauthorized
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, errors.New("internal error")
	}

	return claims, nil
}
