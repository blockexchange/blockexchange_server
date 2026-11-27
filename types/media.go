package types

import (
	"encoding/base64"
	"encoding/json"
)

type Nodedefinition struct {
	Name       string `json:"name" gorm:"primarykey;column:name"`
	ModName    string `json:"mod_name" gorm:"column:mod_name"`
	Definition string `json:"definition" gorm:"column:definition"`
}

func (n *Nodedefinition) TableName() string {
	return "nodedefinition"
}

type Mediafile struct {
	Name string `json:"name" gorm:"primarykey;column:name"`
	// source repository, nil for files from the luanti connections and uploads
	GitMediaSourceID *string `json:"git_media_source_id" gorm:"column:git_media_source_id"`
	Data             []byte  `json:"data" gorm:"column:data"`
}

func (m *Mediafile) TableName() string {
	return "mediafile"
}

func (mf *Mediafile) UnmarshalJSON(data []byte) error {
	m := make(map[string]any)
	err := json.Unmarshal(data, &m)
	if err != nil {
		return err
	}

	mf.Name = getString(m["name"])
	mf.Data, err = base64.RawStdEncoding.DecodeString(getString(m["data"]))
	if err != nil {
		return err
	}

	return nil
}

// mediafile metadata without the binary data
type MediafileInfo struct {
	Name             string  `json:"name" gorm:"column:name"`
	GitMediaSourceID *string `json:"git_media_source_id" gorm:"column:git_media_source_id"`
	// git url and media license of the source, nil without a source
	GitURL       *string `json:"git_url" gorm:"column:git_url"`
	MediaLicense *string `json:"media_license" gorm:"column:media_license"`
	Size         int64   `json:"size" gorm:"column:size"`
}
