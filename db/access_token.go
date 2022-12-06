package db

import (
	"blockexchange/types"
	"database/sql"

	"github.com/minetest-go/dbutil"
)

type AccessTokenRepository struct {
	DB *sql.DB
}

func (r AccessTokenRepository) GetAccessTokensByUserID(user_id int64) ([]*types.AccessToken, error) {
	return dbutil.SelectMulti(r.DB, func() *types.AccessToken { return &types.AccessToken{} }, "where user_id = $1", user_id)
}

func (r AccessTokenRepository) GetAccessTokenByTokenAndUserID(token string, user_id int64) (*types.AccessToken, error) {
	at, err := dbutil.Select(r.DB, &types.AccessToken{}, "where token = $1 and user_id = $2", token, user_id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return at, err
}

func (r AccessTokenRepository) CreateAccessToken(access_token *types.AccessToken) error {
	return dbutil.InsertReturning(r.DB, access_token, "id", &access_token.ID)
}

func (r AccessTokenRepository) IncrementAccessTokenUseCount(id int64) error {
	_, err := r.DB.Exec("update access_token set usecount = usecount + 1 where id = $1", id)
	return err
}

func (r AccessTokenRepository) RemoveAccessToken(id, user_id int64) error {
	_, err := r.DB.Exec("delete from access_token where id = $1 and user_id = $2", id, user_id)
	return err
}
