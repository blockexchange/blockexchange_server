package db

import (
	"blockexchange/types"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SchemaPartRepository struct {
	g *gorm.DB
}

func (r *SchemaPartRepository) CreateOrUpdateSchemaPart(part *types.SchemaPart) error {
	part.OrderID = int64(types.GetSchemaPartOrderID(part.OffsetX, part.OffsetY, part.OffsetZ))

	return r.g.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "schema_uid"}, {Name: "offset_x"}, {Name: "offset_y"}, {Name: "offset_z"}},
		DoUpdates: clause.Assignments(map[string]any{
			"mtime":    part.Mtime,
			"data":     part.Data,
			"metadata": part.MetaData,
		}),
	}).Create(part).Error
}

func (r *SchemaPartRepository) GetBySchemaUIDAndOffset(schema_uid string, offset_x, offset_y, offset_z int) (*types.SchemaPart, error) {
	return FindSingle[types.SchemaPart](r.g.Where(types.SchemaPart{
		SchemaUID: schema_uid,
		OffsetX:   offset_x,
		OffsetY:   offset_y,
		OffsetZ:   offset_z,
	}))
}

func (r *SchemaPartRepository) GetBySchemaUIDAndRange(schema_uid string, x1, y1, z1, x2, y2, z2 int) ([]*types.SchemaPart, error) {
	g := r.g.Model(types.SchemaPart{}).Where("schema_uid = ?", schema_uid)
	g = g.Where("offset_x >= ?", x1).Where("offset_y >= ?", y1).Where("offest_z >= ?", z1)
	g = g.Where("offset_x <= ?", x2).Where("offset_y <= ?", y2).Where("offest_z <= ?", z2)
	return FindMulti[types.SchemaPart](g)
}

func (r *SchemaPartRepository) RemoveBySchemaUIDAndOffset(schema_uid string, offset_x, offset_y, offset_z int) error {
	return r.g.Delete(types.SchemaPart{SchemaUID: schema_uid, OffsetX: offset_x, OffsetY: offset_y, OffsetZ: offset_z}).Error
}

func (r *SchemaPartRepository) GetNextBySchemaUIDAndOffset(schema_uid string, offset_x, offset_y, offset_z int) (*types.SchemaPart, error) {
	order_id := types.GetSchemaPartOrderID(offset_x, offset_y, offset_z)
	g := r.g.Model(types.SchemaPart{}).Where("schema_uid = ?", schema_uid)
	g = g.Where("order_id > ?", order_id)
	g = g.Order("order_id asc")
	return FindSingle[types.SchemaPart](g)
}

func (r *SchemaPartRepository) GetNextBySchemaUIDAndMtime(schema_uid string, mtime int64) (*types.SchemaPart, error) {
	g := r.g.Model(types.SchemaPart{}).Where("schema_uid = ?", schema_uid)
	g = g.Where("mtime > ?", mtime)
	g = g.Order("mtime asc")
	return FindSingle[types.SchemaPart](g)
}

func (r *SchemaPartRepository) CountNextBySchemaUIDAndMtime(schema_uid string, mtime int64) (int64, error) {
	g := r.g.Model(types.SchemaPart{}).Where("schema_uid = ?", schema_uid)
	g = g.Where("mtime > ?", mtime)

	var c int64
	return c, g.Count(&c).Error
}

func (r *SchemaPartRepository) GetFirstBySchemaUID(schema_uid string) (*types.SchemaPart, error) {
	g := r.g.Model(types.SchemaPart{}).Where("schema_uid = ?", schema_uid)
	g = g.Order("mtime asc")
	return FindSingle[types.SchemaPart](g)
}
