package web

import (
	"blockexchange/db"

	"github.com/golang-jwt/jwt/v4"
)

type Context struct {
	JWTKey       string
	CookieName   string
	CookieDomain string
	CookiePath   string
	CookieSecure bool
	Repos        *db.Repositories
}

type Claims struct {
	*jwt.RegisteredClaims
	Username string `json:"username"`
	UserID   int64  `json:"user_id"`
}
