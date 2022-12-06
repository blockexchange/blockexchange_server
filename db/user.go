package db

import (
	"blockexchange/types"
	"database/sql"

	"github.com/minetest-go/dbutil"
)

type UserRepository struct {
	db *sql.DB
}

func (r UserRepository) GetUserById(id int64) (*types.User, error) {
	users, err := dbutil.Select(r.db, &types.User{}, "where id = $1", id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return users, err
}

func (r UserRepository) CountUsers() (int, error) {
	return dbutil.Count(r.db, &types.User{}, "")
}

func (r UserRepository) GetUserByName(name string) (*types.User, error) {
	users, err := dbutil.Select(r.db, &types.User{}, "where name = $1", name)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return users, err
}

func (r UserRepository) GetUserByExternalId(external_id string) (*types.User, error) {
	users, err := dbutil.Select(r.db, &types.User{}, "where external_id = $1", external_id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return users, err
}

func (r UserRepository) GetUsers(limit, offset int) ([]*types.User, error) {
	return dbutil.SelectMulti(r.db, func() *types.User { return &types.User{} }, "limit $1 offset $2", limit, offset)
}

func (r UserRepository) CreateUser(user *types.User) error {
	return dbutil.InsertReturning(r.db, user, "id", &user.ID)
}

func (r UserRepository) UpdateUser(user *types.User) error {
	return dbutil.Update(r.db, user, "where id = $1", user.ID)
}
