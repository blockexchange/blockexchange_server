package db

import (
	"blockexchange/types"
	"context"

	"github.com/google/uuid"
	"github.com/vingarcia/ksql"
)

var tagTable = ksql.NewTable("tag", "uid")

type TagRepository struct {
	kdb ksql.Provider
}

func (r *TagRepository) Create(tag *types.Tag) error {
	if tag.UID == "" {
		tag.UID = uuid.NewString()
	}
	return r.kdb.Insert(context.Background(), tagTable, tag)
}

func (r *TagRepository) Delete(id int64) error {
	return r.kdb.Delete(context.Background(), tagTable, id)
}

func (r *TagRepository) Update(tag *types.Tag) error {
	return r.kdb.Patch(context.Background(), tagTable, tag)
}

func (r *TagRepository) GetByID(id int64) (*types.Tag, error) {
	tag := &types.Tag{}
	err := r.kdb.QueryOne(context.Background(), tag, "from tag where id = $1", id)
	if err == ksql.ErrRecordNotFound {
		return nil, nil
	} else {
		return tag, err
	}
}

func (r *TagRepository) GetAll() ([]*types.Tag, error) {
	list := []*types.Tag{}
	return list, r.kdb.Query(context.Background(), &list, "from tag")
}
