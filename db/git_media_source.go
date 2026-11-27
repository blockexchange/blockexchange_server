package db

import (
	"blockexchange/types"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GitMediaSourceRepository struct {
	g *gorm.DB
}

func (r *GitMediaSourceRepository) GetByID(id string) (*types.GitMediaSource, error) {
	return FindSingle[types.GitMediaSource](r.g.Where(types.GitMediaSource{ID: id}))
}

func (r *GitMediaSourceRepository) GetAll() ([]*types.GitMediaSource, error) {
	return FindMulti[types.GitMediaSource](r.g.Model(types.GitMediaSource{}).Order("git_url"))
}

func (r *GitMediaSourceRepository) Create(s *types.GitMediaSource) error {
	if s.ID == "" {
		s.ID = uuid.NewString()
	}
	return r.g.Create(s).Error
}

func (r *GitMediaSourceRepository) Update(s *types.GitMediaSource) error {
	return r.g.Save(s).Error
}

// removes the source and its media files
func (r *GitMediaSourceRepository) Delete(id string) error {
	return r.g.Delete(types.GitMediaSource{ID: id}).Error
}

// stores the licenses and the contentdb package looked up by the sync, leaves the other fields untouched
func (r *GitMediaSourceRepository) UpdateLicenses(id, cdb_package, media_license, license string) error {
	return r.g.Model(types.GitMediaSource{}).Where("id = ?", id).Updates(map[string]any{
		"cdb_package":   cdb_package,
		"media_license": media_license,
		"license":       license,
	}).Error
}

// writes the result of a sync run, leaves the editable fields untouched
func (r *GitMediaSourceRepository) UpdateStatus(id string, timestamp int64, success bool, message string) error {
	return r.g.Model(types.GitMediaSource{}).Where("id = ?", id).Updates(map[string]any{
		"last_updated_timestamp": timestamp,
		"last_updated_success":   success,
		"last_updated_message":   message,
	}).Error
}
