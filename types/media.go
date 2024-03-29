package types

import (
	"encoding/base64"
	"encoding/json"
)

type Mod struct {
	Name         string `json:"name" ksql:"name"`
	Source       string `json:"source" ksql:"source"`
	MediaLicense string `json:"media_license" ksql:"media_license"`
}

type Nodedefinition struct {
	Name       string `json:"name" ksql:"name"`
	ModName    string `json:"mod_name" ksql:"mod_name"`
	Definition string `json:"definition" ksql:"definition"`
}

type Mediafile struct {
	Name    string `json:"name" ksql:"name"`
	ModName string `json:"mod_name" ksql:"mod_name"`
	Data    []byte `json:"data" ksql:"data"`
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
