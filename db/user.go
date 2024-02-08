package db

import (
	"blockexchange/types"
	"context"

	"github.com/vingarcia/ksql"
)

var userTable = ksql.NewTable("public.user", "id")

type UserRepository struct {
	kdb ksql.Provider
}

func (r *UserRepository) GetUserById(id int64) (*types.User, error) {
	u := &types.User{}
	err := r.kdb.QueryOne(context.Background(), u, "from public.user where id = $1", id)
	if err == ksql.ErrRecordNotFound {
		return nil, nil
	}
	return u, err
}

func (r *UserRepository) CountUsers() (int64, error) {
	c := &types.Count{}
	err := r.kdb.QueryOne(context.Background(), c, "select count(*) as count from public.user")
	return c.Count, err
}

func (r *UserRepository) GetUserByName(name string) (*types.User, error) {
	u := &types.User{}
	err := r.kdb.QueryOne(context.Background(), u, "from public.user where name = $1", name)
	if err == ksql.ErrRecordNotFound {
		return nil, nil
	}
	return u, err
}

func (r *UserRepository) GetUserByExternalIdAndType(external_id string, ut types.UserType) (*types.User, error) {
	u := &types.User{}
	err := r.kdb.QueryOne(context.Background(), u, "from public.user where external_id = $1 and type = $2", external_id, ut)
	if err == ksql.ErrRecordNotFound {
		return nil, nil
	}
	return u, err
}

func (r *UserRepository) GetUsers(limit, offset int) ([]*types.User, error) {
	list := []*types.User{}
	return list, r.kdb.Query(context.Background(), &list, "from public.user limit $1 offset $2", limit, offset)
}

func (r *UserRepository) CreateUser(user *types.User) error {
	return r.kdb.Insert(context.Background(), userTable, user)
}

func (r *UserRepository) UpdateUser(user *types.User) error {
	return r.kdb.Patch(context.Background(), userTable, user)
}
