package testdata

import (
	"blockexchange/types"
)

type MockAccessTokenRepository struct {
}

func (r MockAccessTokenRepository) GetAccessTokensByUserID(user_id int64) ([]types.AccessToken, error) {
	if user_id == 1 {
		return []types.AccessToken{
			types.AccessToken{
				Token:   "abc",
				UserID:  1,
				Expires: 1800,
			},
		}, nil
	} else {
		return []types.AccessToken{}, nil
	}
}
func (r MockAccessTokenRepository) GetAccessTokenByTokenAndUserID(token string, user_id int64) (*types.AccessToken, error) {
	return nil, nil
}
func (r MockAccessTokenRepository) CreateAccessToken(access_token *types.AccessToken) error {
	return nil
}
func (r MockAccessTokenRepository) IncrementAccessTokenUseCount(id int64) error {
	return nil
}
func (r MockAccessTokenRepository) RemoveAccessToken(id, user_id int64) error {
	return nil
}
