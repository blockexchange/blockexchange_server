package db

import (
	"blockexchange/types"
	"context"

	"github.com/google/uuid"
	"github.com/vingarcia/ksql"
)

var accessTokenTable = ksql.NewTable("access_token", "uid")

type AccessTokenRepository struct {
	kdb ksql.Provider
}

func (r *AccessTokenRepository) GetAccessTokensByUserID(user_id int64) ([]*types.AccessToken, error) {
	list := []*types.AccessToken{}
	return list, r.kdb.Query(context.Background(), &list, "from access_token where user_id = $1", user_id)
}

func (r *AccessTokenRepository) GetAccessTokenByTokenAndUserID(token string, user_id int64) (*types.AccessToken, error) {
	t := &types.AccessToken{}
	err := r.kdb.QueryOne(context.Background(), t, "from access_token where token = $1 and user_id = $2", token, user_id)
	if err == ksql.ErrRecordNotFound {
		return nil, nil
	}
	return t, err
}

func (r *AccessTokenRepository) CreateAccessToken(access_token *types.AccessToken) error {
	if access_token.UID == "" {
		access_token.UID = uuid.NewString()
	}
	return r.kdb.Insert(context.Background(), accessTokenTable, access_token)
}

func (r *AccessTokenRepository) IncrementAccessTokenUseCount(uid string) error {
	_, err := r.kdb.Exec(context.Background(), "update access_token set usecount = usecount + 1 where uid = $1", uid)
	return err
}

func (r *AccessTokenRepository) RemoveAccessToken(uid string, user_id int64) error {
	_, err := r.kdb.Exec(context.Background(), "delete from access_token where uid = $1 and user_id = $2", uid, user_id)
	return err
}
