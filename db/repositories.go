package db

import "github.com/jmoiron/sqlx"

type Repositories struct {
	AccessTokenRepo            AccessTokenRepository
	UserRepo                   UserRepository
	SchemaRepo                 SchemaRepository
	SchemaPartRepo             SchemaPartRepository
	SchemaModRepo              SchemaModRepository
	SchemaSearchRepo           SchemaSearchRepository
	SchemaScreenshotRepo       SchemaScreenshotRepository
	CollectionRepo             CollectionRepsoitory
	CollectionSchemaRepository CollectionSchemaRepository
	TagRepo                    TagRepository
	SchemaTagRepo              SchemaTagRepository
	SchemaStarRepo             SchemaStarRepository
}

func NewRepositories(db_ *sqlx.DB) *Repositories {
	return &Repositories{
		AccessTokenRepo:            DBAccessTokenRepository{DB: db_},
		UserRepo:                   DBUserRepository{DB: db_},
		SchemaRepo:                 DBSchemaRepository{DB: db_},
		SchemaPartRepo:             DBSchemaPartRepository{DB: db_},
		SchemaModRepo:              DBSchemaModRepository{DB: db_},
		SchemaSearchRepo:           NewSchemaSearchRepository(db_),
		SchemaScreenshotRepo:       DBSchemaScreenshotRepository{DB: db_},
		CollectionRepo:             DBCollectionRepository{DB: db_},
		CollectionSchemaRepository: DBCollectionSchemaRepository{DB: db_},
		TagRepo:                    DBTagRepository{DB: db_},
		SchemaTagRepo:              DBSchemaTagRepository{DB: db_},
		SchemaStarRepo:             DBSchemaStarRepository{DB: db_},
	}
}
