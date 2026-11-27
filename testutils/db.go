package testutils

import (
	"blockexchange/db"
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func IsDatabaseAvailable() bool {
	return os.Getenv("PGHOST") != ""
}

var (
	testDBOnce  sync.Once
	testDBRepos *db.Repositories
	testDBErr   error
)

// the database of the test process, initialized (connected and migrated) once and shared by all tests
func CreateTestDatabase(t *testing.T) *db.Repositories {
	if !IsDatabaseAvailable() {
		t.SkipNow()
	}

	testDBOnce.Do(func() {
		testDBRepos, testDBErr = db.Init()
	})
	if !assert.NoError(t, testDBErr) {
		t.FailNow()
	}
	assert.NotNil(t, testDBRepos)

	return testDBRepos
}
