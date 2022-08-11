package components

import (
	"blockexchange/db"
	"blockexchange/types"
)

type LatestSchemasModel struct {
	List []*types.SchemaSearchResult
}

func LatestSchemas(repos *db.Repositories) (*LatestSchemasModel, error) {
	m := &LatestSchemasModel{}
	var err error

	m.List, err = repos.SchemaSearchRepo.Search(&types.SchemaSearchRequest{}, 20, 0)
	if err != nil {
		return nil, err
	}

	return m, nil
}
