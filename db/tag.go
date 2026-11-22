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

// updates all fields, including zero values (restricted=false, empty description)
func (r *TagRepository) Update(tag *types.Tag) error {
	return r.g.Save(tag).Error
}

func (r *TagRepository) GetByUID(uid string) (*types.Tag, error) {
	return FindSingle[types.Tag](r.g.Where(types.Tag{UID: uid}))
}

func (r *TagRepository) GetByName(name string) (*types.Tag, error) {
	return FindSingle[types.Tag](r.g.Where(types.Tag{Name: name}))
}

func (r *TagRepository) GetAll() ([]*types.Tag, error) {
	return FindMulti[types.Tag](r.g.Model(types.Tag{}))
}
