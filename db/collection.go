package db

import (
	"blockexchange/types"
	"database/sql"

	"github.com/minetest-go/dbutil"
)

type CollectionRepository struct {
	DB *sql.DB
}

func (repo CollectionRepository) Create(collection *types.Collection) error {
	c := &types.Collection{}
	return dbutil.InsertReturning(repo.DB, c, "id", &c.ID)
}

func (repo CollectionRepository) Delete(id int64) error {
	_, err := repo.DB.Exec("delete from collection where id = $1", id)
	return err
}

func (repo CollectionRepository) Update(collection *types.Collection) error {
	return dbutil.Update(repo.DB, collection, map[string]any{"id": collection.ID})
}

func (repo CollectionRepository) GetByID(id int64) (*types.Collection, error) {
	c, err := dbutil.Select(repo.DB, &types.Collection{}, "id = $1", id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return c, err
}

func (repo CollectionRepository) GetByUserID(user_id int64) ([]*types.Collection, error) {
	return dbutil.SelectMulti(repo.DB, func() *types.Collection { return &types.Collection{} }, "user_id = $1", user_id)
}
