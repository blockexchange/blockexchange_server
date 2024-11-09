package db

import (
	"blockexchange/types"

	"gorm.io/gorm"
)

type SchemaStarRepository struct {
	g *gorm.DB
}

func (r *SchemaStarRepository) Create(ss *types.SchemaStar) error {
	return r.g.Create(ss).Error
}

func (r *SchemaStarRepository) Delete(ss *types.SchemaStar) error {
	return r.g.Delete(ss).Error
}

func (r *SchemaStarRepository) GetBySchemaUID(schema_uid string) ([]*types.SchemaStar, error) {
	return FindMulti[types.SchemaStar](r.g.Where("schema_uid = ?", schema_uid))
}

func (r *SchemaStarRepository) GetBySchemaAndUserID(schema_uid string, user_uid string) (*types.SchemaStar, error) {
	return FindSingle[types.SchemaStar](r.g.Where("schema_uid = ?", schema_uid).Where("user_uid = ?", user_uid))
}

func (r *SchemaStarRepository) CountBySchemaUID(schema_uid string) (int64, error) {
	var c int64
	return c, r.g.Model(types.SchemaStar{}).Where("schema_uid = ?", schema_uid).Count(&c).Error
}

func (r *SchemaStarRepository) CountByUserUID(user_uid string) (int64, error) {
	var c int64
	return c, r.g.Raw("select count(*) from user_schema_star uss join schema s on s.uid = uss.schema_uid where s.user_uid = ?", user_uid).Scan(&c).Error
}
