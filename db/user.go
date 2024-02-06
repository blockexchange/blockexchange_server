package db

import (
	"blockexchange/types"
	"database/sql"

	"github.com/minetest-go/dbutil"
)

func NewUserRepository(DB *sql.DB) *UserRepository {
	return &UserRepository{
		DB:  DB,
		dbu: dbutil.New(DB, dbutil.DialectPostgres, func() *types.User { return &types.User{} }),
	}
}

type UserRepository struct {
	DB  *sql.DB
	dbu *dbutil.DBUtil[*types.User]
}

func (r *UserRepository) GetUserById(id int64) (*types.User, error) {
	users, err := r.dbu.Select("where id = %s", id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return users, err
}

func (r *UserRepository) CountUsers() (int, error) {
	return r.dbu.Count("")
}

func (r *UserRepository) GetUserByName(name string) (*types.User, error) {
	user, err := r.dbu.Select("where name = %s", name)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return user, err
}

func (r *UserRepository) GetUserByExternalIdAndType(external_id string, ut types.UserType) (*types.User, error) {
	user, err := r.dbu.Select("where external_id = %s and type = %s", external_id, ut)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return user, err
}

func (r *UserRepository) GetUsers(limit, offset int) ([]*types.User, error) {
	return r.dbu.SelectMulti("limit %s offset %s", limit, offset)
}

func (r *UserRepository) CreateUser(user *types.User) error {
	return r.dbu.InsertReturning(user, "id", &user.ID)
}

func (r *UserRepository) UpdateUser(user *types.User) error {
	return r.dbu.Update(user, "where id = %s", user.ID)
}
