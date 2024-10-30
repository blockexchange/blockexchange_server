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
		UserRepo:             &UserRepository{kdb: kdb},
		SchemaRepo:           &SchemaRepository{kdb: kdb},
		SchemaPartRepo:       &SchemaPartRepository{kdb: kdb},
		SchemaModRepo:        &SchemaModRepository{kdb: kdb},
		SchemaSearchRepo:     &SchemaSearchRepository{kdb: kdb, DB: DB},
		SchemaScreenshotRepo: &SchemaScreenshotRepository{kdb: kdb},
		TagRepo:              &TagRepository{kdb: kdb},
		CollectionRepo:       &CollectionRepository{g: g},
		SchemaTagRepo:        &SchemaTagRepository{kdb: kdb},
		SchemaStarRepo:       &SchemaStarRepository{kdb: kdb},
		MetaRepository:       &MetaRepository{kdb: kdb},
		MediaRepo:            &MediaRepository{kdb: kdb},
		Lock:                 NewLock(DB),
	}
}
