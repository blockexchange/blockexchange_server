package web

import (
	"blockexchange/db"
	"html/template"

	"github.com/golang-jwt/jwt/v4"
)

type Context struct {
	JWTKey         string
	CookieName     string
	CookieDomain   string
	CookiePath     string
	CookieSecure   bool
	BaseURL        string
	Repos          *db.Repositories
	error_template *template.Template
}

type Claims struct {
	*jwt.RegisteredClaims
	Username string `json:"username"`
	UserID   int64  `json:"user_id"`
}
