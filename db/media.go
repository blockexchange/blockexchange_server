package db

import (
	"blockexchange/types"

	"gorm.io/gorm"
)

type MediaRepository struct {
	g *gorm.DB
}

// mod

func (r *MediaRepository) CreateMod(m *types.Mod) error {
	return r.g.Create(m).Error
}

func (r *MediaRepository) UpdateMod(m *types.Mod) error {
	return r.g.Updates(m).Error
}

func (r *MediaRepository) GetModByName(name string) (*types.Mod, error) {
	return FindSingle[types.Mod](r.g.Where(types.Mod{Name: name}))
}

func (r *MediaRepository) GetMods() ([]*types.Mod, error) {
	return FindMulti[types.Mod](r.g.Model(types.Mod{}))
}

func (r *MediaRepository) RemoveMod(name string) error {
	return r.g.Delete(types.Mod{Name: name}).Error
}

// nodedef

func (r *MediaRepository) CreateNodedefinition(nd *types.Nodedefinition) error {
	return r.g.Create(nd).Error
}

func (r *MediaRepository) UpdateNodedefinition(nd *types.Nodedefinition) error {
	return r.g.Updates(nd).Error
}

func (r *MediaRepository) GetNodedefinitionByName(name string) (*types.Nodedefinition, error) {
	return FindSingle[types.Nodedefinition](r.g.Where(types.Nodedefinition{Name: name}))
}

func (r *MediaRepository) GetNodedefinitions() ([]*types.Nodedefinition, error) {
	return FindMulti[types.Nodedefinition](r.g.Model(types.Nodedefinition{}))
}

func (r *MediaRepository) RemoveNodedefinition(name string) error {
	return r.g.Delete(types.Nodedefinition{Name: name}).Error
}

// mediafile

func (r *MediaRepository) CreateMediafile(mf *types.Mediafile) error {
	return r.g.Create(mf).Error
}

func (r *MediaRepository) UpdateMediafile(mf *types.Mediafile) error {
	return r.g.Updates(mf).Error
}

func (r *MediaRepository) GetMediafileByName(name string) (*types.Mediafile, error) {
	return FindSingle[types.Mediafile](r.g.Where(types.Mediafile{Name: name}))
}

func (r *MediaRepository) GetMediafiles() ([]*types.Mediafile, error) {
	return FindMulti[types.Mediafile](r.g.Model(types.Mediafile{}))
}

func (r *MediaRepository) RemoveMediafile(name string) error {
	return r.g.Delete(types.Mediafile{Name: name}).Error
}
