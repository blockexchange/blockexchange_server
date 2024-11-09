package db

import (
	"blockexchange/types"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TagRepository struct {
	g *gorm.DB
}

func (r *TagRepository) Create(tag *types.Tag) error {
	if tag.UID == "" {
		tag.UID = uuid.NewString()
	}
	return r.g.Create(tag).Error
}

func (r *TagRepository) Update(tag *types.Tag) error {
	return r.g.Updates(tag).Error
}

func (r *TagRepository) GetAll() ([]*types.Tag, error) {
	return FindMulti[types.Tag](r.g.Model(types.Tag{}))
}
