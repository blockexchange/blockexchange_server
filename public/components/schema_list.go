package components

import (
	"blockexchange/controller"
	"blockexchange/types"
)

type SchemaListModel struct {
	Schemas []*types.SchemaSearchResult
	BaseURL string
}

func SchemaList(rc *controller.RenderContext, schemas []*types.SchemaSearchResult) *SchemaListModel {
	return &SchemaListModel{
		Schemas: schemas,
		BaseURL: rc.BaseURL(),
	}
}
