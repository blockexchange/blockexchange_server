package db

import (
	"blockexchange/types"

	"github.com/jmoiron/sqlx"
)

type SchemaSearchRepository interface {
	FindByKeywords(keywords string) ([]types.SchemaSearchResult, error)
	FindRecent(count int) ([]types.SchemaSearchResult, error)
}

func NewSchemaSearchRepository(db *sqlx.DB) SchemaSearchRepository {
	return &DBSchemaSearchRepository{
		DB:       db,
		UserRepo: DBUserRepository{DB: db},
	}
}

type DBSchemaSearchRepository struct {
	DB       *sqlx.DB
	UserRepo UserRepository
}

func (repo DBSchemaSearchRepository) enhance(list []types.Schema) ([]types.SchemaSearchResult, error) {
	result := make([]types.SchemaSearchResult, len(list))

	for i, schema := range list {
		user, err := repo.UserRepo.GetUserById(schema.UserID)
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
		}
	}

	return result, nil
}

func (repo DBSchemaSearchRepository) findSingle(where string, params ...interface{}) (*types.SchemaSearchResult, error) {
	return nil, nil
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
