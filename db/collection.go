package db

import (
	"blockexchange/types"
	"context"

	"github.com/google/uuid"
	"github.com/vingarcia/ksql"
)

var collectionTable = ksql.NewTable("collection", "uid")

type CollectionRepository struct {
	kdb ksql.Provider
}

func (r *CollectionRepository) GetCollectionsByUserUID(user_uid string) ([]*types.Collection, error) {
	list := []*types.Collection{}
	return list, r.kdb.Query(context.Background(), &list, "from collection where user_uid = $1", user_uid)
}

func (r *CollectionRepository) GetCollectionByUserUIDAndName(user_uid, name string) (*types.Collection, error) {
	c := &types.Collection{}
	err := r.kdb.QueryOne(context.Background(), c, "from collection where user_uid = $1 and name = $2", user_uid, name)
	if err == ksql.ErrRecordNotFound {
		return nil, nil
	}
	return c, err
}

func (r *CollectionRepository) CreateCollection(c *types.Collection) error {
	if c.UID == "" {
		c.UID = uuid.NewString()
	}
	return r.kdb.Insert(context.Background(), collectionTable, c)
}

func (r *CollectionRepository) UpdateCollection(c *types.Collection) error {
	return r.kdb.Patch(context.Background(), collectionTable, c)
}

func (r *CollectionRepository) DeleteCollection(collection_uid string) error {
	return r.kdb.Delete(context.Background(), collectionTable, collection_uid)
}
