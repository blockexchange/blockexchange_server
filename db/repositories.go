package db

import (
	"github.com/jmoiron/sqlx"
)

type Repositories struct {
	AccessTokenRepo      *AccessTokenRepository
	UserRepo             *UserRepository
	SchemaRepo           *SchemaRepository
	SchemaPartRepo       *SchemaPartRepository
	SchemaModRepo        *SchemaModRepository
	SchemaSearchRepo     *SchemaSearchRepository
	SchemaScreenshotRepo *SchemaScreenshotRepository
	TagRepo              *TagRepository
	SchemaTagRepo        *SchemaTagRepository
	SchemaStarRepo       *SchemaStarRepository
	MetaRepository       *MetaRepository
}

func NewRepositories(db_ *sqlx.DB) *Repositories {

	return &Repositories{
		AccessTokenRepo:      NewAccessTokenRepository(db_.DB),
		UserRepo:             NewUserRepository(db_.DB),
		SchemaRepo:           NewSchemaRepository(db_.DB),
		SchemaPartRepo:       NewSchemaPartRepository(db_.DB),
		SchemaModRepo:        &SchemaModRepository{DB: db_},
		SchemaSearchRepo:     NewSchemaSearchRepository(db_.DB),
		SchemaScreenshotRepo: &SchemaScreenshotRepository{DB: db_},
		TagRepo:              &TagRepository{DB: db_},
		SchemaTagRepo:        &SchemaTagRepository{DB: db_},
		SchemaStarRepo:       &SchemaStarRepository{DB: db_},
		MetaRepository:       &MetaRepository{db: db_.DB},
	}
}
