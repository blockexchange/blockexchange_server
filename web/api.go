package web

import (
	"blockexchange/core"
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
	CollectionRepo       db.CollectionRepsoitory
	TagRepo              db.TagRepository
	Cache                core.Cache
}

func NewApi(db_ *sqlx.DB, cache core.Cache) *Api {
	return &Api{
		AccessTokenRepo:      db.DBAccessTokenRepository{DB: db_},
		UserRepo:             db.DBUserRepository{DB: db_},
		SchemaRepo:           db.DBSchemaRepository{DB: db_},
		SchemaPartRepo:       db.DBSchemaPartRepository{DB: db_},
		SchemaModRepo:        db.DBSchemaModRepository{DB: db_},
		SchemaSearchRepo:     db.NewSchemaSearchRepository(db_),
		SchemaScreenshotRepo: db.DBSchemaScreenshotRepository{DB: db_},
		CollectionRepo:       db.DBCollectionRepository{DB: db_},
		TagRepo:              db.DBTagRepository{DB: db_},
		Cache:                cache,
	}
}
