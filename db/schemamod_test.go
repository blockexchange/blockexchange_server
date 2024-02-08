package db_test

import (
	"blockexchange/db"
	"blockexchange/types"
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetSchemaModsByIDs(t *testing.T) {
	kdb := CreateTestDatabase(t)
	repos := db.NewRepositories(kdb)

	u := &types.User{
		Name: fmt.Sprintf("test_%d", rand.Intn(1000)),
	}
	assert.NoError(t, repos.UserRepo.CreateUser(u))
	s := &types.Schema{
		UserID: *u.ID,
		Name:   "test",
	}
	assert.NoError(t, repos.SchemaRepo.CreateSchema(s))
	// TODO
}
