package db

import (
	"blockexchange/types"

	"github.com/jmoiron/sqlx"
)

type CollectionSchemaRepository interface {
	Create(collection_id, schema_id int64) error
	Delete(collection_id, schema_id int64) error
	GetBySchemaID(schema_id int64) ([]types.CollectionSchema, error)
	GetByCollectionID(collection_id int64) ([]types.CollectionSchema, error)
}

type DBCollectionSchemaRepository struct {
	DB *sqlx.DB
}

func (repo DBCollectionSchemaRepository) Create(collection_id, schema_id int64) error {
	_, err := repo.DB.Exec("insert into collection_schema(collection_id, schema_id) values($1, $2)", collection_id, schema_id)
	return err
}

func (repo DBCollectionSchemaRepository) Delete(collection_id, schema_id int64) error {
	_, err := repo.DB.Exec("delete from collection_schema where collection_id = $1 and schema_id = $2", collection_id, schema_id)
	return err
}

func (repo DBCollectionSchemaRepository) GetBySchemaID(schema_id int64) ([]types.CollectionSchema, error) {
	list := []types.CollectionSchema{}
	err := repo.DB.Select(&list, "select * from collection_schema where schema_id = $1", schema_id)
	if err != nil {
		return nil, err
	} else {
		return list, nil
	}
}

func (repo DBCollectionSchemaRepository) GetByCollectionID(collection_id int64) ([]types.CollectionSchema, error) {
	list := []types.CollectionSchema{}
	err := repo.DB.Select(&list, "select * from collection_schema where collection_id = $1", collection_id)
	if err != nil {
		return nil, err
	} else {
		return list, nil
	}
}
