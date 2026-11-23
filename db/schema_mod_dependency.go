package db

import (
	"blockexchange/types"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SchemaModDependencyRepository struct {
	g *gorm.DB
}

func (r *SchemaModDependencyRepository) GetByUID(uid string) (*types.SchemaModDependency, error) {
	// explicit conditions: an empty value in a struct condition is ignored by gorm
	return FindSingle[types.SchemaModDependency](r.g.Where("uid = ?", uid))
}

func (r *SchemaModDependencyRepository) GetBySchemaUID(schema_uid string) ([]*types.SchemaModDependency, error) {
	return FindMulti[types.SchemaModDependency](r.g.Where("schema_uid = ?", schema_uid).Order("git_url"))
}

func (r *SchemaModDependencyRepository) Create(d *types.SchemaModDependency) error {
	if d.UID == "" {
		d.UID = uuid.NewString()
	}
	if d.Version == "" {
		d.Version = types.SchemaModDependencyDefaultVersion
	}
	return r.g.Create(d).Error
}

func (r *SchemaModDependencyRepository) Update(d *types.SchemaModDependency) error {
	if d.Version == "" {
		d.Version = types.SchemaModDependencyDefaultVersion
	}
	return r.g.Save(d).Error
}

func (r *SchemaModDependencyRepository) Delete(uid string) error {
	return r.g.Where("uid = ?", uid).Delete(types.SchemaModDependency{}).Error
}
