package db_test

import (
	"blockexchange/db"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func IsDatabaseAvailable() bool {
	return os.Getenv("PGHOST") != ""
}

func CreateTestDatabase(t *testing.T) *db.Repositories {
	if !IsDatabaseAvailable() {
		t.SkipNow()
	}

	repos, err := db.Init()
	assert.NoError(t, err)
	assert.NotNil(t, repos)

	return repos
}
