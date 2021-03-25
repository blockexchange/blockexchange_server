package testutils

import (
	"blockexchange/db"
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func IsDatabaseAvailable() bool {
	return os.Getenv("PGHOST") != ""
}

func CreateTestDatabase(t *testing.T) *sqlx.DB {
	if !IsDatabaseAvailable() {
		t.SkipNow()
	}

	db_, err := db.Init()
	assert.NoError(t, err)
	assert.NotNil(t, db_)
	db.Migrate(db_.DB)

	return db_
}
