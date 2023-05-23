package web

import (
	"blockexchange/core"
	"blockexchange/db"
)

type Context struct {
	JWTKey       string
	CookieName   string
	CookieDomain string
	CookiePath   string
	CookieSecure bool
	BaseURL      string
	Repos        *db.Repositories
	Config       *core.Config
}
