package db

import (
	"blockexchange/types"

	"gorm.io/gorm"
)

type SchemaPullClientRepository struct {
	g *gorm.DB
}

func (r *SchemaPullClientRepository) Create(sp *types.SchematicPullClient) error {
	return r.g.Create(sp).Error
}

func (r *SchemaPullClientRepository) GetBySchemaUID(schema_uid string) ([]*types.SchematicPullClient, error) {
	return FindMulti[types.SchematicPullClient](r.g.Where("schema_pull_uid = ?", schema_uid))
}

func (r *SchemaPullClientRepository) Update(sp *types.SchematicPullClient) error {
	return r.g.Select("*").Updates(sp).Error
}

func (r *SchemaPullClientRepository) Delete(uid string) error {
	return r.g.Where("uid = ?", uid).Delete(types.SchematicPullClient{}).Error
}
