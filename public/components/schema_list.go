package components

import (
	"blockexchange/controller"
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

func SchemaList(rc *controller.RenderContext, schemas []*types.SchemaSearchResult, showuser bool) *SchemaListModel {
	m := &SchemaListModel{
		Schemas:  make([]*SchemaListEntry, len(schemas)),
		ShowUser: showuser,
		BaseURL:  rc.BaseURL(),
	}

	for i, s := range schemas {
		m.Schemas[i] = &SchemaListEntry{
			Schema: s,
		}

		if rc.Claims() != nil && rc.Claims().Username == s.UserName {
			m.Schemas[i].Owner = true
		}
	}
	return m
}
