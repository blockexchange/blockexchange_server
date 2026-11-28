package db

import (
	"blockexchange/types"
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SchemaPartRepository struct {
	g *gorm.DB
}

// UpdateWithNextMtime runs fn in a transaction with the next mtime of the schema.
//
// The schema mtime is set to max(now, mtime + 1), so every change gets an mtime that is unique per schema
// and the by-mtime paging (mtime > ?) can't skip parts with the same mtime. The update locks the schema row
// until the transaction is committed: concurrent changes to the same schema are serialized and become visible
// in mtime order (a reader never sees a newer mtime before an older one is committed).
func (r *SchemaPartRepository) UpdateWithNextMtime(schema_uid string, fn func(repo *SchemaPartRepository, mtime int64) error) error {
	return r.g.Transaction(func(tx *gorm.DB) error {
		var mtime int64
		err := tx.Raw("update schema set mtime = greatest(mtime + 1, ?) where uid = ? returning mtime",
			time.Now().UnixMilli(), schema_uid).Scan(&mtime).Error
		if err != nil {
			return err
		}
		if mtime == 0 {
			return fmt.Errorf("schema not found: %s", schema_uid)
		}
		return fn(&SchemaPartRepository{g: tx}, mtime)
	})
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
	g := r.g.Model(types.SchemaPart{})
	g = g.Where("offset_x = ?", offset_x)
	g = g.Where("offset_y = ?", offset_y)
	g = g.Where("offset_z = ?", offset_z)
	g = g.Where("schema_uid = ?", schema_uid)
	return FindSingle[types.SchemaPart](g)
}

func (r *SchemaPartRepository) GetBySchemaUIDAndRange(schema_uid string, x1, y1, z1, x2, y2, z2 int) ([]*types.SchemaPart, error) {
	g := r.g.Model(types.SchemaPart{})
	g = g.Where("schema_uid = ?", schema_uid)
	g = g.Where("offset_x >= ?", x1).Where("offset_y >= ?", y1).Where("offset_z >= ?", z1)
	g = g.Where("offset_x <= ?", x2).Where("offset_y <= ?", y2).Where("offset_z <= ?", z2)
	return FindMulti[types.SchemaPart](g)
}

func (r *SchemaPartRepository) GetNextBySchemaUIDAndOffset(schema_uid string, offset_x, offset_y, offset_z int) (*types.SchemaPart, error) {
	order_id := types.GetSchemaPartOrderID(offset_x, offset_y, offset_z)
	g := r.g.Model(types.SchemaPart{})
	g = g.Where("schema_uid = ?", schema_uid)
	g = g.Where("order_id > ?", order_id)
	g = g.Order("order_id asc")
	return FindSingle[types.SchemaPart](g)
}

func (r *SchemaPartRepository) GetNextBySchemaUIDAndMtime(schema_uid string, mtime int64) (*types.SchemaPart, error) {
	g := r.g.Model(types.SchemaPart{})
	g = g.Where("schema_uid = ?", schema_uid)
	g = g.Where("mtime > ?", mtime)
	// mtimes are unique per schema since UpdateWithNextMtime, order_id keeps older duplicates deterministic
	g = g.Order("mtime asc").Order("order_id asc")
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
	g = g.Order("order_id asc")
	return FindSingle[types.SchemaPart](g)
}
