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

func (r *AccessTokenRepository) GetAccessTokensByUserUID(user_uid string) ([]*types.AccessToken, error) {
	list := []*types.AccessToken{}
	return list, r.kdb.Query(context.Background(), &list, "from access_token where user_uid = $1", user_uid)
}

func (r *AccessTokenRepository) GetAccessTokenByTokenAndUserUID(token string, user_uid string) (*types.AccessToken, error) {
	t := &types.AccessToken{}
	err := r.kdb.QueryOne(context.Background(), t, "from access_token where token = $1 and user_uid = $2", token, user_uid)
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

func (r *AccessTokenRepository) RemoveAccessToken(uid string, user_uid string) error {
	_, err := r.kdb.Exec(context.Background(), "delete from access_token where uid = $1 and user_uid = $2", uid, user_uid)
	return err
}
