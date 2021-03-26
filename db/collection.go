package db

import (
	"blockexchange/types"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type CollectionRepsoitory interface {
	Create(collection *types.Collection) error
	Delete(id int64) error
	Update(collection *types.Collection) error
	GetByID(id int64) (*types.Collection, error)
	GetByUserID(user_id int64) ([]types.Collection, error)
}

type DBCollectionRepository struct {
	DB *sqlx.DB
}

func (repo DBCollectionRepository) Create(collection *types.Collection) error {
	logrus.Trace("db.CreateCollection", collection)
	query := `
		insert into
		collection(user_id, name, description)
		values(:user_id, :name, :description)
		returning id
	`
	stmt, err := repo.DB.PrepareNamed(query)
	if err != nil {
		return err
	}
	return stmt.Get(&collection.ID, collection)
}

func (repo DBCollectionRepository) Delete(id int64) error {
	_, err := repo.DB.Exec("delete from collection where id = $1", id)
	return err
}

func (repo DBCollectionRepository) Update(collection *types.Collection) error {
	query := `
	update collection
	set
		name = :name,
		description = :description
	where id = :id
`
	_, err := repo.DB.NamedExec(query, collection)
	return err
}

func (repo DBCollectionRepository) GetByID(id int64) (*types.Collection, error) {
	collection := types.Collection{}
	err := repo.DB.Get(&collection, "select * from collection where id = $1", id)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	} else {
		return &collection, nil
	}
}

func (repo DBCollectionRepository) GetByUserID(user_id int64) ([]types.Collection, error) {
	list := []types.Collection{}
	query := `select * from collection where user_id = $1`
	err := repo.DB.Select(&list, query, user_id)
	if err != nil {
		return nil, err
	} else {
		return list, nil
	}
}
