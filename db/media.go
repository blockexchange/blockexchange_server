package db

import (
	"blockexchange/types"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// MediaQuery filters and pages node definitions and media files
type MediaQuery struct {
	// case-insensitive substring match on the name or mod name, empty matches all
	Query  string
	Limit  int
	Offset int
}

var likeEscaper = strings.NewReplacer(`\`, `\\`, "%", `\%`, "_", `\_`)

func applyMediaFilter(g *gorm.DB, q string) *gorm.DB {
	if q == "" {
		return g
	}
	pattern := "%" + likeEscaper.Replace(q) + "%"
	return g.Where("name ilike ? or mod_name ilike ?", pattern, pattern)
}

type MediaRepository struct {
	g *gorm.DB
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

func (r *MediaRepository) SearchNodedefinitions(q MediaQuery) ([]*types.Nodedefinition, error) {
	g := applyMediaFilter(r.g.Model(types.Nodedefinition{}), q.Query)
	return FindMulti[types.Nodedefinition](g.Order("name").Limit(q.Limit).Offset(q.Offset))
}

func (r *MediaRepository) CountNodedefinitions(q string) (int64, error) {
	var c int64
	return c, applyMediaFilter(r.g.Model(types.Nodedefinition{}), q).Count(&c).Error
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

func (r *MediaRepository) SearchMediafileInfos(q MediaQuery) ([]*types.MediafileInfo, error) {
	list := make([]*types.MediafileInfo, 0)
	g := applyMediaFilter(r.g.Model(types.Mediafile{}), q.Query)
	err := g.Select("name, mod_name, length(data) as size").Order("name").Limit(q.Limit).Offset(q.Offset).Scan(&list).Error
	return list, err
}

func (r *MediaRepository) CountMediafiles(q string) (int64, error) {
	var c int64
	return c, applyMediaFilter(r.g.Model(types.Mediafile{}), q).Count(&c).Error
}

// inserts or replaces the given node definitions (by name)
func (r *MediaRepository) UpsertNodedefinitions(list []*types.Nodedefinition) error {
	if len(list) == 0 {
		return nil
	}
	return r.g.Clauses(clause.OnConflict{UpdateAll: true}).CreateInBatches(list, 500).Error
}

// inserts or replaces the given media files (by name)
func (r *MediaRepository) UpsertMediafiles(list []*types.Mediafile) error {
	if len(list) == 0 {
		return nil
	}
	return r.g.Clauses(clause.OnConflict{UpdateAll: true}).CreateInBatches(list, 50).Error
}
