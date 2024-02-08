package testutils

import (
	"blockexchange/db"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vingarcia/ksql"
)

func IsDatabaseAvailable() bool {
	return os.Getenv("PGHOST") != ""
}

func CreateTestDatabase(t *testing.T) ksql.Provider {
	if !IsDatabaseAvailable() {
		t.SkipNow()
	}

	kdb, err := db.Init()
	assert.NoError(t, err)
	assert.NotNil(t, kdb)

	return kdb
}
