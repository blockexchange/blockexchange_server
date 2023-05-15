package web

import "github.com/golang-jwt/jwt/v4"

type Context struct {
	JWTKey       string
	CookieName   string
	CookieDomain string
	CookiePath   string
	CookieSecure bool
}

type Claims struct {
	*jwt.RegisteredClaims
	Username string `json:"username"`
	UserID   int64  `json:"user_id"`
}
