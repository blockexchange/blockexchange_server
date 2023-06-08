package components

import (
	"blockexchange/types"
)

type SchemaListModel struct {
	Schemas  []*SchemaListEntry
	BaseURL  string
	ShowUser bool
}

type SchemaListEntry struct {
	Owner  bool
	Schema *types.SchemaSearchResult
}

func SchemaList(c *types.Claims, schemas []*types.SchemaSearchResult, showuser bool) *SchemaListModel {
	m := &SchemaListModel{
		Schemas:  make([]*SchemaListEntry, len(schemas)),
		ShowUser: showuser,
	}

	for i, s := range schemas {
		m.Schemas[i] = &SchemaListEntry{
			Schema: s,
		}

		if c != nil && c.Username == s.UserName {
			m.Schemas[i].Owner = true
		}
	}
	return m
}
