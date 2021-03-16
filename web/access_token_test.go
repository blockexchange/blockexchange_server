package web

import (
	"blockexchange/types"
	"blockexchange/web/testdata"
	"encoding/json"
	"net/http/httptest"
	"testing"
)

func TestGetAccessTokensNoToken(t *testing.T) {
	r := httptest.NewRequest("GET", "http://", nil)
	w := httptest.NewRecorder()
	tokenInfo := types.TokenInfo{}
	ctx := SecureContext{Token: &tokenInfo}

	repo := testdata.MockAccessTokenRepository{}
	api := Api{
		AccessTokenRepo: &repo,
	}

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

	repo := testdata.MockAccessTokenRepository{}
	api := Api{
		AccessTokenRepo: &repo,
	}

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

	repo := testdata.MockAccessTokenRepository{}
	api := Api{
		AccessTokenRepo: &repo,
	}

	api.GetAccessTokens(w, r, &ctx)
	result := []types.AccessToken{}
	json.NewDecoder(w.Body).Decode(&result)
	if len(result) != 1 {
		t.Fatal("invalid length")
	}
}
