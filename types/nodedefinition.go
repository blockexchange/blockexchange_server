package types

type Nodedefinition struct {
	Name       string `json:"name" ksql:"name"`
	ModName    string `json:"mod_name" ksql:"mod_name"`
	Definition string `json:"definition" ksql:"definition"`
}
