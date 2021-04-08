package db

import (
	"blockexchange/types"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type SchemaSearchRepository interface {
	FindByKeywords(keywords string) ([]types.SchemaSearchResult, error)
	FindRecent(count int) ([]types.SchemaSearchResult, error)
	FindByUsernameAndSchemaname(schema_name, user_name string) (*types.SchemaSearchResult, error)
}

func NewSchemaSearchRepository(db *sqlx.DB) SchemaSearchRepository {
	return &DBSchemaSearchRepository{
		DB:            db,
		UserRepo:      DBUserRepository{DB: db},
		SchemaModRepo: DBSchemaModRepository{DB: db},
		SchemaTagRepo: DBSchemaTagRepository{DB: db},
	}
}

type DBSchemaSearchRepository struct {
	DB            *sqlx.DB
	UserRepo      UserRepository
	SchemaModRepo SchemaModRepository
	SchemaTagRepo SchemaTagRepository
}

func (repo DBSchemaSearchRepository) enhance(list []types.Schema) ([]types.SchemaSearchResult, error) {
	result := make([]types.SchemaSearchResult, len(list))

	for i, schema := range list {
		user, err := repo.UserRepo.GetUserById(schema.UserID)
		if err != nil {
			return nil, err
		}

		mods, err := repo.SchemaModRepo.GetSchemaModsBySchemaID(schema.ID)
		if err != nil {
			return nil, err
		}

		mod_list := make([]string, len(mods))
		for i, mod := range mods {
			mod_list[i] = mod.ModName
		}

		tag_list, err := repo.SchemaTagRepo.GetBySchemaID(schema.ID)
		if err != nil {
			return nil, err
		}

		result[i] = types.SchemaSearchResult{
			Schema: schema,
			User: &types.User{
				Name:    user.Name,
				ID:      user.ID,
				Created: user.Created,
				Type:    user.Type,
			},
			Mods: mod_list,
			Tags: tag_list,
		}
	}

	return result, nil
}

func (repo DBSchemaSearchRepository) findSingle(where string, params ...interface{}) (*types.SchemaSearchResult, error) {
	list, err := repo.findMulti(where, params...)
	if err != nil {
		return nil, err
	} else if len(list) == 1 {
		return &list[0], nil
	} else {
		return nil, nil
	}
}

func (repo DBSchemaSearchRepository) findMulti(where string, params ...interface{}) ([]types.SchemaSearchResult, error) {
	list := []types.Schema{}
	err := repo.DB.Select(&list, "select * from schema where "+where, params...)
	if err != nil {
		return nil, err
	}

	return repo.enhance(list)
}

func (repo DBSchemaSearchRepository) FindByKeywords(keywords string) ([]types.SchemaSearchResult, error) {
	return repo.findMulti("search_tokens @@ to_tsquery($1)", keywords)
}

func (repo DBSchemaSearchRepository) FindRecent(count int) ([]types.SchemaSearchResult, error) {
	return repo.findMulti("complete = true order by created desc limit $1", count)
}

func (repo DBSchemaSearchRepository) FindByUsernameAndSchemaname(schema_name, user_name string) (*types.SchemaSearchResult, error) {
	logrus.WithFields(logrus.Fields{
		"schema_name": schema_name,
		"user_name":   user_name,
	}).Trace("DBSchemaSearchRepository::FindByUsernameAndSchemaname")

	where := `name = $1 and user_id = (select id from public.user where name = $2)`
	return repo.findSingle(where, schema_name, user_name)
}
