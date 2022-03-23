package types

type ColorResponse struct {
	Red   uint8 `json:"r"`
	Green uint8 `json:"g"`
	Blue  uint8 `json:"b"`
	Alpha uint8 `json:"a"`
}

type ColorMappingResponse struct {
	Colors map[string]*ColorResponse `json:"colors"`
}
