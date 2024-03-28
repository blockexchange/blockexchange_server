package db

import (
	"blockexchange/types"
	"context"

	"github.com/vingarcia/ksql"
)

var modTable = ksql.NewTable("mod", "name")

type MediaRepository struct {
	kdb ksql.Provider
}

func (r *MediaRepository) CreateMod(m *types.Mod) error {
	return r.kdb.Insert(context.Background(), modTable, m)
}

func (r *MediaRepository) UpdateMod(m *types.Mod) error {
	return r.kdb.Patch(context.Background(), modTable, m)
}

func (r *MediaRepository) GetModByName(name string) (*types.Mod, error) {
	m := &types.Mod{}
	err := r.kdb.QueryOne(context.Background(), m, "from mod where name = $1", name)
	if err == ksql.ErrRecordNotFound {
		return nil, nil
	} else {
		return m, err
	}
}

func (r *MediaRepository) GetMods() ([]*types.Mod, error) {
	list := []*types.Mod{}
	err := r.kdb.Query(context.Background(), &list, "from mod")
	return list, err
}

func (r *MediaRepository) RemoveMod(name string) error {
	return r.kdb.Delete(context.Background(), modTable, name)
}
