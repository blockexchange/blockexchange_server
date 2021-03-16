package web

import "blockexchange/db"

type Api struct {
	AccessTokenRepo db.AccessTokenRepository
	UserRepo        db.UserRepository
	SchemaRepo      db.SchemaRepository
}
