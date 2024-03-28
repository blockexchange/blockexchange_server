package types

type Mod struct {
	Name         string `json:"name" ksql:"name"`
	Source       string `json:"source" ksql:"source"`
	CodeLicense  string `json:"code_license" ksql:"code_license"`
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
