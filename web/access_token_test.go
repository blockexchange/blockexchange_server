package web

import (
	"blockexchange/types"
	"encoding/json"
	"net/http/httptest"
	"testing"
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

func TestGetAccessTokensNoToken(t *testing.T) {
	r := httptest.NewRequest("GET", "http://", nil)
	w := httptest.NewRecorder()
	tokenInfo := types.TokenInfo{}
	ctx := SecureContext{Token: &tokenInfo}

	repo := MockAccessTokenRepository{}
	api := AccessTokenApi{Repo: &repo}

	api.GetAccessTokens(w, r, &ctx)
}

func TestGetAccessTokensInvalidUser(t *testing.T) {
	r := httptest.NewRequest("GET", "http://", nil)
	w := httptest.NewRecorder()
	tokenInfo := types.TokenInfo{
		UserID: 2,
		Permissions: []types.JWTPermission{
			types.JWTPermissionManagement,
		},
	}
	ctx := SecureContext{Token: &tokenInfo}

	repo := MockAccessTokenRepository{}
	api := AccessTokenApi{Repo: &repo}

	api.GetAccessTokens(w, r, &ctx)
	result := []types.AccessToken{}
	json.NewDecoder(w.Body).Decode(&result)
	if len(result) != 0 {
		t.Fatal("invalid length")
	}
}

func TestGetAccessTokensValidUser(t *testing.T) {
	r := httptest.NewRequest("GET", "http://", nil)
	w := httptest.NewRecorder()
	tokenInfo := types.TokenInfo{
		UserID: 1,
		Permissions: []types.JWTPermission{
			types.JWTPermissionManagement,
		},
	}
	ctx := SecureContext{Token: &tokenInfo}

	repo := MockAccessTokenRepository{}
	api := AccessTokenApi{Repo: &repo}

	api.GetAccessTokens(w, r, &ctx)
	result := []types.AccessToken{}
	json.NewDecoder(w.Body).Decode(&result)
	if len(result) != 1 {
		t.Fatal("invalid length")
	}
}
