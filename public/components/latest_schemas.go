package components

import (
	"blockexchange/db"
	"blockexchange/types"
)

type LatestSchemasModel struct {
	Schemas []*SchemaCardModel
}

func LatestSchemas(baseUrl string, repos *db.Repositories) (*LatestSchemasModel, error) {
	m := &LatestSchemasModel{}
	var err error

	schemas, err := repos.SchemaSearchRepo.Search(&types.SchemaSearchRequest{}, 12, 0)
	if err != nil {
		return nil, err
	}

	m.Schemas = make([]*SchemaCardModel, len(schemas))
	for i, s := range schemas {
		m.Schemas[i] = SchemaCard(baseUrl, s)
	}

	return m, nil
}
