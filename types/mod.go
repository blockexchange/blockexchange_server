package types

type Mod struct {
	Name         string `json:"name" ksql:"name"`
	Source       string `json:"source" ksql:"source"`
	CodeLicense  string `json:"code_license" ksql:"code_license"`
	MediaLicense string `json:"media_license" ksql:"media_license"`
}
