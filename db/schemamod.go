package db

import (
	"blockexchange/types"

	"gorm.io/gorm"
)

type SchemaModRepository struct {
	g *gorm.DB
}

func (r *SchemaModRepository) GetSchemaModsBySchemaUID(schema_uid string) ([]*types.SchemaMod, error) {
	return FindMulti[types.SchemaMod](r.g.Where(types.SchemaMod{SchemaUID: schema_uid}))
}

func (r *SchemaModRepository) GetSchemaModCount() ([]*types.SchemaModCount, error) {
	return FindMulti[types.SchemaModCount](r.g.Select("mod_name", "count(*) as count").Table("schemamod").Group("mod_name").Order("count desc"))
}

func (r *SchemaModRepository) CreateSchemaMod(schema_mod *types.SchemaMod) error {
	return r.g.Create(schema_mod).Error
}

func (r *SchemaModRepository) RemoveSchemaMods(schema_uid string) error {
	return r.g.Where("schema_uid = ?", schema_uid).Delete(types.SchemaMod{}).Error
}
