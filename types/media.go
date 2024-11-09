package types

import (
	"encoding/base64"
	"encoding/json"
)

type Mod struct {
	Name         string `json:"name" gorm:"primarykey;column:name"`
	Source       string `json:"source" gorm:"column:source"`
	MediaLicense string `json:"media_license" gorm:"column:media_license"`
}

func (m *Mod) TableName() string {
	return "mod"
}

type Nodedefinition struct {
	Name       string `json:"name" gorm:"primarykey;column:name"`
	ModName    string `json:"mod_name" gorm:"column:mod_name"`
	Definition string `json:"definition" gorm:"column:definition"`
}

func (n *Nodedefinition) TableName() string {
	return "nodedefinition"
}

type Mediafile struct {
	Name    string `json:"name" gorm:"primarykey;column:name"`
	ModName string `json:"mod_name" gorm:"column:mod_name"`
	Data    []byte `json:"data" gorm:"column:data"`
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

	mf.ModName = getString(m["mod_name"])
	mf.Name = getString(m["name"])
	mf.Data, err = base64.RawStdEncoding.DecodeString(getString(m["data"]))
	if err != nil {
		return err
	}

	return nil
}
