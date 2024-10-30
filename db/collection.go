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
	list := []*types.Collection{}
	err := r.g.Where(types.Collection{UID: uid}).Limit(1).Find(&list).Error
	if err != nil || len(list) == 0 {
		return nil, err
	} else {
		return list[0], nil
	}
}

func (r *CollectionRepository) GetCollectionsByUserUID(user_uid string) ([]*types.Collection, error) {
	list := []*types.Collection{}
	err := r.g.Where(types.Collection{UserUID: user_uid}).Find(&list).Error
	return list, err
}

func (r *CollectionRepository) GetCollectionByUserUIDAndName(user_uid, name string) (*types.Collection, error) {
	list := []*types.Collection{}
	err := r.g.Where(types.Collection{UserUID: user_uid, Name: name}).Limit(1).Find(&list).Error
	if err != nil || len(list) == 0 {
		return nil, err
	} else {
		return list[0], nil
	}
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
