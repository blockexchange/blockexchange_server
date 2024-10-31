package db

import (
	"blockexchange/types"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CollectionRepository struct {
	g *gorm.DB
}

func (r *CollectionRepository) GetCollectionByUID(uid string) (*types.Collection, error) {
	return FindSingle[types.Collection](r.g.Where(types.Collection{UID: uid}))
}

func (r *CollectionRepository) GetCollectionsByUserUID(user_uid string) ([]*types.Collection, error) {
	return FindMulti[types.Collection](r.g.Where(types.Collection{UserUID: user_uid}))
}

func (r *CollectionRepository) GetCollectionByUserUIDAndName(user_uid, name string) (*types.Collection, error) {
	return FindSingle[types.Collection](r.g.Where(types.Collection{UserUID: user_uid, Name: name}))
}

func (r *CollectionRepository) CreateCollection(c *types.Collection) error {
	if c.UID == "" {
		c.UID = uuid.NewString()
	}
	return r.g.Create(c).Error
}

func (r *CollectionRepository) UpdateCollection(c *types.Collection) error {
	return r.g.Save(c).Error
}

func (r *CollectionRepository) DeleteCollection(collection_uid string) error {
	return r.g.Delete(types.Collection{UID: collection_uid}).Error
}
