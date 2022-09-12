package db

import (
	"blockexchange/types"
	"database/sql"
)

type CollectionSchemaRepository struct {
	DB *sql.DB
}

func (repo CollectionSchemaRepository) Create(collection_id, schema_id int64) error {
	return Insert(repo.DB, &types.CollectionSchema{CollectionID: collection_id, SchemaID: schema_id})
}

func (repo CollectionSchemaRepository) Delete(collection_id, schema_id int64) error {
	_, err := repo.DB.Exec("delete from collection_schema where collection_id = $1 and schema_id = $2", collection_id, schema_id)
	return err
}

func (repo CollectionSchemaRepository) GetBySchemaID(schema_id int64) ([]*types.CollectionSchema, error) {
	return SelectMulti(repo.DB, func() *types.CollectionSchema { return &types.CollectionSchema{} }, "where schema_id = $1", schema_id)
}

func (repo CollectionSchemaRepository) GetByCollectionID(collection_id int64) ([]*types.CollectionSchema, error) {
	return SelectMulti(repo.DB, func() *types.CollectionSchema { return &types.CollectionSchema{} }, "where collection_id = $1", collection_id)
}
