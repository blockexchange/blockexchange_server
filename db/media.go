package db

import (
	"blockexchange/types"
	"slices"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// MediaQuery filters and pages node definitions and media files
type MediaQuery struct {
	// case-insensitive substring match on the name or mod name, empty matches all
	Query string
	// media files only: true matches files with a media license, false files without one, nil all
	License *bool
	Limit   int
	Offset  int
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

// returns the node definitions matching the given names, unknown names are skipped
func (r *MediaRepository) GetNodedefinitionsByNames(names []string) ([]*types.Nodedefinition, error) {
	if len(names) == 0 {
		return []*types.Nodedefinition{}, nil
	}
	return FindMulti[types.Nodedefinition](r.g.Where("name in ?", names).Order("name"))
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

// media files have no mod name, the filter matches the file name only
// joins the git source as "s" for the license filter
func applyMediafileFilter(g *gorm.DB, q MediaQuery) *gorm.DB {
	g = g.Joins("left join git_media_source s on s.id = mediafile.git_media_source_id")
	if q.Query != "" {
		g = g.Where("mediafile.name ilike ?", "%"+likeEscaper.Replace(q.Query)+"%")
	}
	if q.License != nil {
		if *q.License {
			g = g.Where("coalesce(s.media_license, '') <> ''")
		} else {
			g = g.Where("coalesce(s.media_license, '') = ''")
		}
	}
	return g
}

func (r *MediaRepository) CreateMediafile(mf *types.Mediafile) error {
	return r.g.Create(mf).Error
}

func (r *MediaRepository) UpdateMediafile(mf *types.Mediafile) error {
	return r.g.Updates(mf).Error
}

func (r *MediaRepository) GetMediafileByName(name string) (*types.Mediafile, error) {
	return FindSingle[types.Mediafile](r.g.Where(types.Mediafile{Name: name}))
}

// returns the media files matching the given names, unknown names are skipped
func (r *MediaRepository) GetMediafilesByNames(names []string) ([]*types.Mediafile, error) {
	if len(names) == 0 {
		return []*types.Mediafile{}, nil
	}
	return FindMulti[types.Mediafile](r.g.Where("name in ?", names).Order("name"))
}

func (r *MediaRepository) GetMediafiles() ([]*types.Mediafile, error) {
	return FindMulti[types.Mediafile](r.g.Model(types.Mediafile{}))
}

func (r *MediaRepository) RemoveMediafile(name string) error {
	return r.g.Delete(types.Mediafile{Name: name}).Error
}

func (r *MediaRepository) SearchMediafileInfos(q MediaQuery) ([]*types.MediafileInfo, error) {
	list := make([]*types.MediafileInfo, 0)
	g := applyMediafileFilter(r.g.Model(types.Mediafile{}), q).
		Select("mediafile.name, mediafile.git_media_source_id, s.git_url, s.media_license, length(mediafile.data) as size")
	err := g.Order("mediafile.name").Limit(q.Limit).Offset(q.Offset).Scan(&list).Error
	return list, err
}

func (r *MediaRepository) CountMediafiles(q MediaQuery) (int64, error) {
	var c int64
	return c, applyMediafileFilter(r.g.Model(types.Mediafile{}), q).Count(&c).Error
}

// names of the media files imported from the git source
func (r *MediaRepository) GetMediafileNamesBySource(source_id string) ([]string, error) {
	names := []string{}
	err := r.g.Model(types.Mediafile{}).Where("git_media_source_id = ?", source_id).Order("name").Pluck("name", &names).Error
	return names, err
}

// removes the media files of the git source with the given names
func (r *MediaRepository) RemoveMediafilesOfSource(source_id string, names []string) error {
	for chunk := range slices.Chunk(names, 1000) {
		err := r.g.Where("git_media_source_id = ? and name in ?", source_id, chunk).Delete(&types.Mediafile{}).Error
		if err != nil {
			return err
		}
	}
	return nil
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
