package db

import (
	"blockexchange/types"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AccessTokenRepository struct {
	g *gorm.DB
}

func (r *AccessTokenRepository) GetAccessTokensByUserUID(user_uid string) ([]*types.AccessToken, error) {
	list := []*types.AccessToken{}
	err := r.g.Where(types.AccessToken{UserUID: user_uid}).Find(&list).Error
	return list, err
}

func (r *AccessTokenRepository) GetAccessTokenByTokenAndUserUID(token string, user_uid string) (*types.AccessToken, error) {
	list := []*types.AccessToken{}
	err := r.g.Where(types.AccessToken{UserUID: user_uid, Token: token}).Limit(1).Find(&list).Error
	if err != nil || len(list) == 0 {
		return nil, err
	} else {
		return list[0], nil
	}
}

func (r *AccessTokenRepository) CreateAccessToken(access_token *types.AccessToken) error {
	if access_token.UID == "" {
		access_token.UID = uuid.NewString()
	}
	return r.g.Create(access_token).Error
}

func (r *AccessTokenRepository) IncrementAccessTokenUseCount(uid string) error {
	return r.g.Exec("update access_token set usecount = usecount + 1 where uid = ?", uid).Error
}

func (r *AccessTokenRepository) RemoveAccessToken(uid string, user_uid string) error {
	return r.g.Delete(types.AccessToken{UID: uid, UserUID: user_uid}).Error
}
