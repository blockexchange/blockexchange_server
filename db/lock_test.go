package db_test

import (
	"testing"
	"time"

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

	err = lock.RunLocked(id, 1*time.Second, func() error {
		return nil
	})
	assert.NoError(t, err)
}
