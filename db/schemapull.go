package db

import (
	"blockexchange/types"

	"gorm.io/gorm"
)

type SchemaPullRepository struct {
	g *gorm.DB
}

func (r *SchemaPullRepository) Create(sp *types.SchematicPull) error {
	return r.g.Create(sp).Error
}

func (r *SchemaPullRepository) GetBySchemaUID(schema_uid string) (*types.SchematicPull, error) {
	return FindSingle[types.SchematicPull](r.g.Where("schema_uid = ?", schema_uid))
}

func (r *SchemaPullRepository) GetByNextRun(next_run_lt int64) ([]*types.SchematicPull, error) {
	return FindMulti[types.SchematicPull](r.g.Where("enabled = ?", true).Where("next_run < ?", next_run_lt))
}

func (r *SchemaPullRepository) Update(sp *types.SchematicPull) error {
	return r.g.Select("*").Updates(sp).Error
}

func (r *SchemaPullRepository) Delete(schema_uid string) error {
	return r.g.Where("schema_uid = ?", schema_uid).Delete(types.SchematicPull{}).Error
}