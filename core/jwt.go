package core

import (
	"blockexchange/types"
	"errors"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var err_unauthorized = errors.New("unauthorized")

func GetPermissions(user *types.User, management bool) []types.JWTPermission {
	permissions := []types.JWTPermission{
		types.JWTPermissionUpload,
		types.JWTPermissionOverwrite,
	}

	if management {
		permissions = append(permissions, types.JWTPermissionManagement)
	}

	if user.Role == types.UserRoleAdmin && management {
		permissions = append(permissions, types.JWTPermissionAdmin)
	}

	return permissions
}

func CreateClaims(user *types.User, permissions []types.JWTPermission) *types.Claims {
	return &types.Claims{
		UserID:      *user.ID,
		Username:    user.Name,
		Type:        user.Type,
		Permissions: permissions,
	}
}

func (c *Core) SetClaims(w http.ResponseWriter, token string, d time.Duration) {
	co := &http.Cookie{
		Name:     c.cfg.CookieName,
		Value:    token,
		Path:     c.cfg.CookiePath,
		Expires:  time.Now().Add(d),
		HttpOnly: true,
		Secure:   c.cfg.CookieSecure,
	}
	http.SetCookie(w, co)
}

func (c *Core) RemoveClaims(w http.ResponseWriter) {
	co := &http.Cookie{
		Name:     c.cfg.CookieName,
		Value:    "",
		Path:     c.cfg.CookiePath,
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   c.cfg.CookieSecure,
	}
	http.SetCookie(w, co)
}

func (c *Core) GetClaims(r *http.Request) (*types.Claims, error) {
	var token string
	authorization := r.Header.Get("Authorization")
	if authorization != "" {
		// token in header
		token = authorization
	} else {
		// token in cookie
		co, err := r.Cookie(c.cfg.CookieName)
		if err == http.ErrNoCookie {
			return nil, nil
		}
		if err != nil {
			return nil, err
		}
		token = co.Value
	}
	if token == "" {
		// no token found
		return nil, nil
	}

	return c.ParseJWT(token)
}

func (c *Core) CreateJWT(user *types.User, permissions []types.JWTPermission, d time.Duration) (string, error) {
	cl := CreateClaims(user, permissions)

	cl.RegisteredClaims = &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(d)),
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, cl)

	return t.SignedString([]byte(c.cfg.Key))

}

func (c *Core) ParseJWT(token string) (*types.Claims, error) {
	t, err := jwt.ParseWithClaims(token, &types.Claims{}, func(token *jwt.Token) (any, error) {
		return []byte(c.cfg.Key), nil
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
