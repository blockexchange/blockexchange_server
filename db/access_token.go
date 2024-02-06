package db

import (
	"blockexchange/types"
	"database/sql"

	"github.com/minetest-go/dbutil"
)

func NewAccessTokenRepository(DB *sql.DB) *AccessTokenRepository {
	return &AccessTokenRepository{
		DB:  DB,
		dbu: dbutil.New(DB, dbutil.DialectPostgres, func() *types.AccessToken { return &types.AccessToken{} }),
	}
}

type AccessTokenRepository struct {
	DB  *sql.DB
	dbu *dbutil.DBUtil[*types.AccessToken]
}

func (r AccessTokenRepository) GetAccessTokensByUserID(user_id int64) ([]*types.AccessToken, error) {
	return r.dbu.SelectMulti("where user_id = %s", user_id)
}

func (r AccessTokenRepository) GetAccessTokenByTokenAndUserID(token string, user_id int64) (*types.AccessToken, error) {
	at, err := r.dbu.Select("where token = %s and user_id = %s", token, user_id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return at, err
}

func (r AccessTokenRepository) CreateAccessToken(access_token *types.AccessToken) error {
	return r.dbu.InsertReturning(access_token, "id", &access_token.ID)
}

func (r AccessTokenRepository) IncrementAccessTokenUseCount(id int64) error {
	_, err := r.DB.Exec("update access_token set usecount = usecount + 1 where id = $1", id)
	return err
}

func (r AccessTokenRepository) RemoveAccessToken(id, user_id int64) error {
	_, err := r.DB.Exec("delete from access_token where id = $1 and user_id = $2", id, user_id)
	return err
}
