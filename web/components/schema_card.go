package components

import "blockexchange/types"

type SchemaCardModel struct {
	Schema *types.SchemaSearchResult
}

func SchemaCard(s *types.SchemaSearchResult) *SchemaCardModel {
	return &SchemaCardModel{
		Schema: s,
	}
}
