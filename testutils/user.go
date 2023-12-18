package testutils

import (
	"blockexchange/db"
	"blockexchange/types"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func CreateName(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func CreateUser(repo *db.UserRepository, t *testing.T, user *types.User) *types.User {
	if user == nil {
		// create new user
		user = &types.User{}
	}
	// set defaults
	if user.Name == "" {
		user.Name = CreateName(10)
	}
	if user.Type == "" {
		user.Type = types.UserTypeLocal
	}

	assert.NoError(t, repo.CreateUser(user))

	return user
}

func CreateAccessToken(repo *db.AccessTokenRepository, t *testing.T, token *types.AccessToken) *types.AccessToken {
	if token == nil {
		token = &types.AccessToken{}
	}

	assert.NoError(t, repo.CreateAccessToken(token))
	return token
}
