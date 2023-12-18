package db

import "github.com/jmoiron/sqlx"

type Repositories struct {
	AccessTokenRepo            *AccessTokenRepository
	UserRepo                   *UserRepository
	SchemaRepo                 *SchemaRepository
	SchemaPartRepo             *SchemaPartRepository
	SchemaModRepo              SchemaModRepository
	SchemaSearchRepo           *SchemaSearchRepository
	SchemaScreenshotRepo       SchemaScreenshotRepository
	CollectionRepo             *CollectionRepository
	CollectionSchemaRepository *CollectionSchemaRepository
	TagRepo                    TagRepository
	SchemaTagRepo              SchemaTagRepository
	SchemaStarRepo             SchemaStarRepository
}

func NewRepositories(db_ *sqlx.DB) *Repositories {
	return &Repositories{
		AccessTokenRepo:            &AccessTokenRepository{DB: db_.DB},
		UserRepo:                   &UserRepository{db: db_.DB},
		SchemaRepo:                 &SchemaRepository{DB: db_.DB},
		SchemaPartRepo:             &SchemaPartRepository{DB: db_.DB},
		SchemaModRepo:              &DBSchemaModRepository{DB: db_},
		SchemaSearchRepo:           &SchemaSearchRepository{DB: db_.DB},
		SchemaScreenshotRepo:       DBSchemaScreenshotRepository{DB: db_},
		CollectionRepo:             &CollectionRepository{DB: db_.DB},
		CollectionSchemaRepository: &CollectionSchemaRepository{DB: db_.DB},
		TagRepo:                    &DBTagRepository{DB: db_},
		SchemaTagRepo:              &DBSchemaTagRepository{DB: db_},
		SchemaStarRepo:             &DBSchemaStarRepository{DB: db_},
	}
}
