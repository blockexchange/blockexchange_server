package db_test

import (
	"blockexchange/types"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestAccessToken(t *testing.T) {
	repos := CreateTestDatabase(t)

	u := &types.User{
		Name: fmt.Sprintf("test_%s", uuid.NewString()),
	}
	assert.NoError(t, repos.UserRepo.CreateUser(u))

	// create

	at := &types.AccessToken{
		UID:      uuid.NewString(),
		Name:     "default",
		Token:    "abcd",
		UserUID:  u.UID,
		Created:  time.Now().Unix(),
		Expires:  time.Now().Add(24 * time.Hour).Unix(),
		UseCount: 0,
	}
	assert.NoError(t, repos.AccessTokenRepo.CreateAccessToken(at))

	// read

	list, err := repos.AccessTokenRepo.GetAccessTokensByUserUID(u.UID)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(list))

	// increment

	assert.NoError(t, repos.AccessTokenRepo.IncrementAccessTokenUseCount(at.UID))

	// delete

	assert.NoError(t, repos.AccessTokenRepo.RemoveAccessToken(at.UID, u.UID))
}
