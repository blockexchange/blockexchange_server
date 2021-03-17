package render

import (
	"embed"
	"encoding/json"
)

//go:embed colormapping.json
var fs embed.FS

type Color struct {
	Red   int `json:"r"`
	Green int `json:"g"`
	Blue  int `json:"b"`
}

func GetColorMapping() (map[string]*Color, error) {
	file, err := fs.Open("colormapping.json")
	if err != nil {
		return nil, err
	}

	m := make(map[string]*Color)
	err = json.NewDecoder(file).Decode(&m)
	return m, err
}
