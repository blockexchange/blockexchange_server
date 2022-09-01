package components

import (
	"blockexchange/controller"
	"blockexchange/types"
)

type SchemaListEntry struct {
	BaseURL string
	Owner   bool
	Schema  *types.SchemaSearchResult
}

func SchemaList(rc *controller.RenderContext, schemas []*types.SchemaSearchResult) []*SchemaListEntry {
	list := make([]*SchemaListEntry, len(schemas))
	for i, s := range schemas {
		list[i] = &SchemaListEntry{
			Schema:  s,
			BaseURL: rc.BaseURL(),
		}

		if rc.Claims() != nil && rc.Claims().Username == s.UserName {
			list[i].Owner = true
		}
	}
	return list
}
