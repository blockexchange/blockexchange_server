package db

import (
	"blockexchange/types"

	"github.com/jmoiron/sqlx"
)

type AccessTokenRepository interface {
	GetAccessTokensByUserID(user_id int64) ([]types.AccessToken, error)
	GetAccessTokenByTokenAndUserID(token string, user_id int64) (*types.AccessToken, error)
	CreateAccessToken(access_token *types.AccessToken) error
	IncrementAccessTokenUseCount(id int64) error
	RemoveAccessToken(id, user_id int64) error
}

type DBAccessTokenRepository struct {
	DB *sqlx.DB
}

func (r DBAccessTokenRepository) GetAccessTokensByUserID(user_id int64) ([]types.AccessToken, error) {
	list := []types.AccessToken{}
	err := r.DB.Select(&list, "select * from access_token where user_id = $1", user_id)
	if err != nil {
		return nil, err
	} else {
		return list, nil
	}
}

func (r DBAccessTokenRepository) GetAccessTokenByTokenAndUserID(token string, user_id int64) (*types.AccessToken, error) {
	list := []types.AccessToken{}
	err := r.DB.Select(&list, "select * from access_token where token = $1 and user_id = $2", token, user_id)
	if err != nil {
		return nil, err
	} else if len(list) == 1 {
		return &list[0], nil
	} else {
		return nil, nil
	}
}

func (r DBAccessTokenRepository) CreateAccessToken(access_token *types.AccessToken) error {
	query := `
		insert into
		access_token(
			name, token, user_id,
			created, expires
		)
		values(
			:name, :token, :user_id,
			:created, :expires
		)
		returning id
	`
	stmt, err := r.DB.PrepareNamed(query)
	if err != nil {
		return err
	}
	return stmt.Get(&access_token.ID, access_token)
}

func (r DBAccessTokenRepository) IncrementAccessTokenUseCount(id int64) error {
	_, err := r.DB.Exec("update access_token set usecount = usecount + 1 where id = $1", id)
	return err
}

func (r DBAccessTokenRepository) RemoveAccessToken(id, user_id int64) error {
	_, err := r.DB.Exec("delete from access_token where id = $1 and user_id = $2", id, user_id)
	return err
}
