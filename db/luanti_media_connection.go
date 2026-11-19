package db

import (
	"blockexchange/types"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LuantiMediaConnectionRepository struct {
	g *gorm.DB
}

func (r *LuantiMediaConnectionRepository) GetByID(id string) (*types.LuantiMediaConnection, error) {
	return FindSingle[types.LuantiMediaConnection](r.g.Where(types.LuantiMediaConnection{ID: id}))
}

func (r *LuantiMediaConnectionRepository) GetAll() ([]*types.LuantiMediaConnection, error) {
	return FindMulti[types.LuantiMediaConnection](r.g.Model(types.LuantiMediaConnection{}))
}

func (r *LuantiMediaConnectionRepository) Create(c *types.LuantiMediaConnection) error {
	if c.ID == "" {
		c.ID = uuid.NewString()
	}
	return r.g.Create(c).Error
}

func (r *LuantiMediaConnectionRepository) Update(c *types.LuantiMediaConnection) error {
	return r.g.Save(c).Error
}

func (r *LuantiMediaConnectionRepository) Delete(id string) error {
	return r.g.Delete(types.LuantiMediaConnection{ID: id}).Error
}
