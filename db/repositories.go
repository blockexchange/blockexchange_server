package db

import (
	"database/sql"

	"github.com/vingarcia/ksql"
	"gorm.io/gorm"
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
	CollectionRepo       *CollectionRepository
	SchemaTagRepo        *SchemaTagRepository
	SchemaStarRepo       *SchemaStarRepository
	MetaRepository       *MetaRepository
	MediaRepo            *MediaRepository
	Lock                 *DBLock
}

func NewRepositories(kdb ksql.Provider, g *gorm.DB, DB *sql.DB) *Repositories {

	return &Repositories{
		AccessTokenRepo:      &AccessTokenRepository{g: g},
		UserRepo:             &UserRepository{g: g},
		SchemaRepo:           &SchemaRepository{g: g},
		SchemaPartRepo:       &SchemaPartRepository{g: g},
		SchemaModRepo:        &SchemaModRepository{g: g},
		SchemaSearchRepo:     &SchemaSearchRepository{kdb: kdb, DB: DB},
		SchemaScreenshotRepo: &SchemaScreenshotRepository{g: g},
		TagRepo:              &TagRepository{g: g},
		CollectionRepo:       &CollectionRepository{g: g},
		SchemaTagRepo:        &SchemaTagRepository{g: g},
		SchemaStarRepo:       &SchemaStarRepository{g: g},
		MetaRepository:       &MetaRepository{g: g},
		MediaRepo:            &MediaRepository{g: g},
		Lock:                 NewLock(DB),
	}
}
