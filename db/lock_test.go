package db_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDBLock(t *testing.T) {
	repos := CreateTestDatabase(t)
	lock := repos.Lock

	id := int64(123)

	result, err := lock.TryLock(id)
	assert.NoError(t, err)
	assert.True(t, result)

	result, err = lock.UnLock(id)
	assert.NoError(t, err)
	assert.True(t, result)

	result, err = lock.UnLock(id)
	assert.NoError(t, err)
	assert.False(t, result)
}
