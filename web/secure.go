package web

import (
	"blockexchange/types"
	"errors"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type SecureHandlerFunc func(http.ResponseWriter, *http.Request, *types.Claims)
type ClaimsCheck func(*types.Claims) (bool, error)

func permissionCheck(req_perms ...types.JWTPermission) ClaimsCheck {
	return func(c *types.Claims) (bool, error) {
		if len(req_perms) > 0 && c == nil {
			return false, errors.New("no credentials found")
		}
		for _, req_perm := range req_perms {
			if !c.HasPermission(req_perm) {
				return false, errors.New("forbidden")
			}
		}
		return true, nil
	}
}

var err_unauthorized = errors.New("unauthorized")
var err_forbidden = errors.New("forbidden")

func (ctx *Context) OptionalSecure(h SecureHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, err := ctx.GetClaims(r)
		if err != nil {
			ctx.RenderError(w, r, 500, err)
			return
		}

		h(w, r, claims)
	}
}

func (ctx *Context) Secure(h SecureHandlerFunc, checks ...ClaimsCheck) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, err := ctx.GetClaims(r)
		if err != nil {
			ctx.RenderError(w, r, 500, err)
			return
		}
		if claims == nil {
			ctx.RenderError(w, r, 401, err_unauthorized)
			return
		}

		for _, c := range checks {
			ok, err := c(claims)
			if err != nil {
				ctx.RenderError(w, r, 500, err)
				return
			}
			if !ok {
				ctx.RenderError(w, r, 403, err_forbidden)
				return
			}
		}

		h(w, r, claims)
	}
}

func (ctx *Context) CreateJWT(c *types.Claims, d time.Duration) (string, error) {
	c.RegisteredClaims = &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(d)),
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return t.SignedString([]byte(ctx.JWTKey))
}

func (ctx *Context) ParseJWT(token string) (*types.Claims, error) {
	t, err := jwt.ParseWithClaims(token, &types.Claims{}, func(token *jwt.Token) (any, error) {
		return []byte(ctx.JWTKey), nil
	})

	if err != nil {
		return nil, err
	}

	if !t.Valid {
		return nil, err_unauthorized
	}

	claims, ok := t.Claims.(*types.Claims)
	if !ok {
		return nil, errors.New("internal error")
	}

	return claims, nil
}

func (ctx *Context) GetClaims(r *http.Request) (*types.Claims, error) {
	co, err := r.Cookie(ctx.CookieName)
	if err == http.ErrNoCookie {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	c, err := ctx.ParseJWT(co.Value)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (ctx *Context) SetClaims(w http.ResponseWriter, token string, d time.Duration) {
	c := &http.Cookie{
		Name:     ctx.CookieName,
		Value:    token,
		Path:     ctx.CookiePath,
		Expires:  time.Now().Add(d),
		Domain:   ctx.CookieDomain,
		HttpOnly: true,
		Secure:   ctx.CookieSecure,
	}
	http.SetCookie(w, c)
}

func (ctx *Context) ClearClaims(w http.ResponseWriter) {
	ctx.SetClaims(w, "", time.Duration(0))
}
