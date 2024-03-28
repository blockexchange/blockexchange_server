package db

import (
	"blockexchange/types"
	"context"

	"github.com/vingarcia/ksql"
)

type MediaRepository struct {
	kdb ksql.Provider
}

// mod

var modTable = ksql.NewTable("mod", "name")

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

// nodedef

var nodedefTable = ksql.NewTable("nodedefinition", "name")

func (r *MediaRepository) CreateNodedefinition(nd *types.Nodedefinition) error {
	return r.kdb.Insert(context.Background(), nodedefTable, nd)
}

func (r *MediaRepository) UpdateNodedefinition(nd *types.Nodedefinition) error {
	return r.kdb.Patch(context.Background(), nodedefTable, nd)
}

func (r *MediaRepository) GetNodedefinitionByName(name string) (*types.Nodedefinition, error) {
	nd := &types.Nodedefinition{}
	err := r.kdb.QueryOne(context.Background(), nd, "from nodedefinition where name = $1", name)
	if err == ksql.ErrRecordNotFound {
		return nil, nil
	} else {
		return nd, err
	}
}
func (r *MediaRepository) GetNodedefinitions() ([]*types.Nodedefinition, error) {
	list := []*types.Nodedefinition{}
	err := r.kdb.Query(context.Background(), &list, "from nodedefinition")
	return list, err
}

func (r *MediaRepository) RemoveNodedefinition(name string) error {
	return r.kdb.Delete(context.Background(), nodedefTable, name)
}

// mediafile

var mediafileTable = ksql.NewTable("mediafile", "name")

func (r *MediaRepository) CreateMediafile(mf *types.Mediafile) error {
	return r.kdb.Insert(context.Background(), mediafileTable, mf)
}

func (r *MediaRepository) UpdateMediafile(mf *types.Mediafile) error {
	return r.kdb.Patch(context.Background(), mediafileTable, mf)
}

func (r *MediaRepository) GetMediafileByName(name string) (*types.Mediafile, error) {
	mf := &types.Mediafile{}
	err := r.kdb.QueryOne(context.Background(), mf, "from mediafile where name = $1", name)
	if err == ksql.ErrRecordNotFound {
		return nil, nil
	} else {
		return mf, err
	}
}

func (r *MediaRepository) GetMediafiles() ([]*types.Mediafile, error) {
	list := []*types.Mediafile{}
	err := r.kdb.Query(context.Background(), &list, "from mediafile")
	return list, err
}

func (r *MediaRepository) RemoveMediafile(name string) error {
	return r.kdb.Delete(context.Background(), mediafileTable, name)
}
