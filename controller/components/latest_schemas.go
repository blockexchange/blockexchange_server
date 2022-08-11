package components

import (
	"blockexchange/db"
)

type LatestSchemasModel struct {
}

func LatestSchemas(repos *db.Repositories) (*LatestSchemasModel, error) {
	m := &LatestSchemasModel{}
	return m, nil
}
