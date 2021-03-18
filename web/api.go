package web

import (
	"blockexchange/db"

	"github.com/jmoiron/sqlx"
)

type Api struct {
	AccessTokenRepo      db.AccessTokenRepository
	UserRepo             db.UserRepository
	SchemaRepo           db.SchemaRepository
	SchemaPartRepo       db.SchemaPartRepository
	SchemaModRepo        db.SchemaModRepository
	SchemaSearchRepo     db.SchemaSearchRepository
	SchemaScreenshotRepo db.SchemaScreenshotRepository
}

func NewApi(db_ *sqlx.DB) *Api {
	return &Api{
		AccessTokenRepo:      db.DBAccessTokenRepository{DB: db_},
		UserRepo:             db.DBUserRepository{DB: db_},
		SchemaRepo:           db.DBSchemaRepository{DB: db_},
		SchemaPartRepo:       db.DBSchemaPartRepository{DB: db_},
		SchemaModRepo:        db.DBSchemaModRepository{DB: db_},
		SchemaSearchRepo:     db.NewSchemaSearchRepository(db_),
		SchemaScreenshotRepo: db.DBSchemaScreenshotRepository{DB: db_},
	}
}
