package types

type GitMediaSource struct {
	ID                   string `json:"id" gorm:"primarykey;column:id"`
	GitURL               string `json:"git_url" gorm:"column:git_url"`
	CDBPackage           string `json:"cdb_package" gorm:"column:cdb_package"`
	MediaLicense         string `json:"media_license" gorm:"column:media_license"`
	License              string `json:"license" gorm:"column:license"`
	LastUpdatedTimestamp int64  `json:"last_updated_timestamp" gorm:"column:last_updated_timestamp"`
	LastUpdatedSuccess   bool   `json:"last_updated_success" gorm:"column:last_updated_success"`
	LastUpdatedMessage   string `json:"last_updated_message" gorm:"column:last_updated_message"`
}

func (s *GitMediaSource) TableName() string {
	return "git_media_source"
}
