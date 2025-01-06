package db

import (
	"blockexchange/types"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SchemaRepository struct {
	g *gorm.DB
}

func (r *SchemaRepository) GetSchemaByUID(uid string) (*types.Schema, error) {
	return FindSingle[types.Schema](r.g.Where("uid = ?", uid))
}

func (r *SchemaRepository) GetSchemaByUserUIDAndName(user_uid string, name string) (*types.Schema, error) {
	return FindSingle[types.Schema](r.g.Where("user_uid = ?", user_uid).Where("name = ?", name))
}

func (r *SchemaRepository) GetSchemaByUsernameAndName(username, schemaname string) (*types.Schema, error) {
	return FindSingle[types.Schema](r.g.Where("user_uid = (select uid from public.user where name = ?)", username).Where("name = ?", schemaname))
}

func (r *SchemaRepository) CreateSchema(schema *types.Schema) error {
	if schema.UID == "" {
		schema.UID = uuid.NewString()
	}
	return r.g.Create(schema).Error
}

func (r *SchemaRepository) UpdateSchema(schema *types.Schema) error {
	return r.g.Updates(schema).Error
}

func (r *SchemaRepository) IncrementViews(uid string) error {
	return r.g.Exec("update schema set views = views + 1 where uid = ?", uid).Error
}

func (r *SchemaRepository) IncrementDownloads(uid string) error {
	return r.g.Exec("update schema set downloads = downloads + 1 where uid = ?", uid).Error
}

func (r *SchemaRepository) DeleteSchema(uid string) error {
	return r.g.Delete(types.Schema{UID: uid}).Error
}

func (r *SchemaRepository) DeleteIncompleteSchema(user_uid string, name string) error {
	return r.g.Where("user_uid = ?", user_uid).Where("name = ?", name).Where("complete = ?", false).Delete(types.Schema{}).Error
}

func (r *SchemaRepository) DeleteOldIncompleteSchema(time_before int64) error {
	return r.g.Where("created < ?", time_before).Where("complete = ?", false).Delete(types.Schema{}).Error
}

func (r *SchemaRepository) CalculateStats(uid string) error {
	return r.g.Exec(`
		update schema s
		set total_size = (
			select
			coalesce(sum(length(data)) + sum(length(metadata)), 0)
			from schemapart sp where sp.schema_uid = s.uid
		),
		total_parts = (select count(*) from schemapart sp where sp.schema_uid = s.uid),
		stars = (select count(*) from user_schema_star where schema_uid = s.uid)
		where s.uid = ?;
	`, uid).Error
}

func (r *SchemaRepository) GetTotalSize() (int64, error) {
	var c int64
	return c, r.g.Raw("select coalesce(sum(total_size),0) as count from schema").Scan(&c).Error
}
