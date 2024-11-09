package db

import (
	"blockexchange/types"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SchemaTagRepository struct {
	g *gorm.DB
}

func (r *SchemaTagRepository) Create(st *types.SchemaTag) error {
	if st.UID == "" {
		st.UID = uuid.NewString()
	}
	return r.g.Create(st).Error
}

func (r *SchemaTagRepository) Delete(schema_uid string, tag_uid string) error {
	return r.g.Delete(types.SchemaTag{SchemaUID: schema_uid, TagUID: tag_uid}).Error
}

func (r *SchemaTagRepository) GetBySchemaUID(schema_uid string) ([]*types.SchemaTag, error) {
	return FindMulti[types.SchemaTag](r.g.Where(types.SchemaTag{SchemaUID: schema_uid}))
}
