package db

import (
	"blockexchange/types"
	"sync"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type SchemaSearchRepository interface {
	FindByKeywords(keywords string, limit, offset int) ([]types.SchemaSearchResult, error)
	FindRecent(limit, offset int) ([]types.SchemaSearchResult, error)
	FindByTagID(tag_id int64) ([]types.SchemaSearchResult, error)
	FindByUserID(tag_id int64) ([]types.SchemaSearchResult, error)
	FindBySchemaID(schema_id int64) ([]types.SchemaSearchResult, error)
	FindByUsername(user_name string) ([]types.SchemaSearchResult, error)
	FindByUsernameAndSchemaname(schema_name, user_name string) (*types.SchemaSearchResult, error)
}

func NewSchemaSearchRepository(db *sqlx.DB) SchemaSearchRepository {
	return &DBSchemaSearchRepository{
		DB:             db,
		UserRepo:       DBUserRepository{DB: db},
		SchemaModRepo:  DBSchemaModRepository{DB: db},
		SchemaTagRepo:  DBSchemaTagRepository{DB: db},
		SchemaStarRepo: DBSchemaStarRepository{DB: db},
	}
}

type DBSchemaSearchRepository struct {
	DB             *sqlx.DB
	UserRepo       UserRepository
	SchemaModRepo  SchemaModRepository
	SchemaTagRepo  SchemaTagRepository
	SchemaStarRepo SchemaStarRepository
}

func (repo DBSchemaSearchRepository) enhance_schema(schema *types.Schema) (*types.SchemaSearchResult, error) {
	user, err := repo.UserRepo.GetUserById(schema.UserID)
	if err != nil {
		return nil, err
	}

	stars, err := repo.SchemaStarRepo.CountBySchemaID(schema.ID)
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

	result := &types.SchemaSearchResult{
		Schema: *schema,
		Stars:  stars,
		User: &types.User{
			Name:    user.Name,
			ID:      user.ID,
			Created: user.Created,
			Type:    user.Type,
		},
		Mods: mod_list,
		Tags: tag_list,
	}

	return result, nil
}

func (repo DBSchemaSearchRepository) enhance(list []*types.Schema) ([]types.SchemaSearchResult, error) {
	result := make([]types.SchemaSearchResult, len(list))
	var wg sync.WaitGroup
	err_chan := make(chan error)

	for i, schema := range list {
		wg.Add(1)
		go func(schema *types.Schema, index int, wg *sync.WaitGroup, err_chan chan error) {
			enhanced_schema, err := repo.enhance_schema(schema)
			if err != nil {
				err_chan <- err
			} else {
				result[index] = *enhanced_schema
			}
			wg.Done()
		}(schema, i, &wg, err_chan)
	}

	wg.Wait()

	select {
	case err := <-err_chan:
		return nil, err
	default:
		return result, nil
	}
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
	list := []*types.Schema{}
	err := repo.DB.Select(&list, "select * from schema where "+where, params...)
	if err != nil {
		return nil, err
	}

	return repo.enhance(list)
}

func (repo DBSchemaSearchRepository) FindByKeywords(keywords string, limit, offset int) ([]types.SchemaSearchResult, error) {
	return repo.findMulti("search_tokens @@ to_tsquery($1) limit $2 offset $3", keywords, limit, offset)
}

func (repo DBSchemaSearchRepository) FindRecent(limit, offset int) ([]types.SchemaSearchResult, error) {
	return repo.findMulti("complete = true order by created desc limit $1 offset $2", limit, offset)
}

func (repo DBSchemaSearchRepository) FindByUsernameAndSchemaname(schema_name, user_name string) (*types.SchemaSearchResult, error) {
	logrus.WithFields(logrus.Fields{
		"schema_name": schema_name,
		"user_name":   user_name,
	}).Trace("DBSchemaSearchRepository::FindByUsernameAndSchemaname")

	where := `name = $1 and user_id = (select id from public.user where name = $2)`
	return repo.findSingle(where, schema_name, user_name)
}

func (repo DBSchemaSearchRepository) FindByUsername(user_name string) ([]types.SchemaSearchResult, error) {
	where := `user_id = (select id from public.user where name = $1)`
	return repo.findMulti(where, user_name)
}

func (repo DBSchemaSearchRepository) FindByTagID(tag_id int64) ([]types.SchemaSearchResult, error) {
	return repo.findMulti("complete = true and id in (select schema_id from schematag where tag_id = $1)", tag_id)
}

func (repo DBSchemaSearchRepository) FindByUserID(tag_id int64) ([]types.SchemaSearchResult, error) {
	return repo.findMulti("complete = true and user_id = $1", tag_id)
}

func (repo DBSchemaSearchRepository) FindBySchemaID(schema_id int64) ([]types.SchemaSearchResult, error) {
	return repo.findMulti("complete = true and id = $1", schema_id)
}
