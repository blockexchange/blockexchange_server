package components

import "blockexchange/types"

type SchemaCardModel struct {
	BaseURL string
	Schema  *types.SchemaSearchResult
}

func SchemaCard(baseUrl string, s *types.SchemaSearchResult) *SchemaCardModel {
	return &SchemaCardModel{
		BaseURL: baseUrl,
		Schema:  s,
	}
}
