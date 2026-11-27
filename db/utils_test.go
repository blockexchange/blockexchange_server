package db_test

import (
	"blockexchange/db"
	"blockexchange/testutils"
	"testing"
)

func CreateTestDatabase(t *testing.T) *db.Repositories {
	return testutils.CreateTestDatabase(t)
}
