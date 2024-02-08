package db

import (
	"github.com/vingarcia/ksql"
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

func NewRepositories(kdb ksql.Provider) *Repositories {

	return &Repositories{
		AccessTokenRepo:      &AccessTokenRepository{kdb: kdb},
		UserRepo:             &UserRepository{kdb: kdb},
		SchemaRepo:           &SchemaRepository{kdb: kdb},
		SchemaPartRepo:       &SchemaPartRepository{kdb: kdb},
		SchemaModRepo:        &SchemaModRepository{kdb: kdb},
		SchemaSearchRepo:     &SchemaSearchRepository{kdb: kdb},
		SchemaScreenshotRepo: &SchemaScreenshotRepository{kdb: kdb},
		TagRepo:              &TagRepository{kdb: kdb},
		SchemaTagRepo:        &SchemaTagRepository{kdb: kdb},
		SchemaStarRepo:       &SchemaStarRepository{kdb: kdb},
		MetaRepository:       &MetaRepository{kdb: kdb},
	}
}
