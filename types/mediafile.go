package types

type Mediafile struct {
	Name    string `json:"name" ksql:"name"`
	ModName string `json:"mod_name" ksql:"mod_name"`
	Data    []byte `json:"data" ksql:"data"`
}
