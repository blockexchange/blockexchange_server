package db_test

import (
	"blockexchange/db"
	"blockexchange/testutils"
	"blockexchange/types"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func createPart(schema *types.Schema, x int) *types.SchemaPart {
	return &types.SchemaPart{
		SchemaUID: schema.UID,
		OffsetX:   x,
		Data:      []byte{0x01},
		MetaData:  []byte{0x02},
	}
}

func TestUpdateWithNextMtime(t *testing.T) {
	repos := CreateTestDatabase(t)
	user := testutils.CreateUser(repos.UserRepo, t, nil)
	schema := testutils.CreateSchema(repos.SchemaRepo, t, user, nil)

	save := func(x int) int64 {
		var mtime int64
		err := repos.SchemaPartRepo.UpdateWithNextMtime(schema.UID, func(repo *db.SchemaPartRepository, m int64) error {
			mtime = m
			part := createPart(schema, x)
			part.Mtime = m
			return repo.CreateOrUpdateSchemaPart(part)
		})
		assert.NoError(t, err)
		return mtime
	}

	// mtimes are unique and increasing, even within the same millisecond
	before := time.Now().UnixMilli()
	m1 := save(0)
	m2 := save(16)
	m3 := save(0)
	assert.GreaterOrEqual(t, m1, before)
	assert.Greater(t, m2, m1)
	assert.Greater(t, m3, m2)

	stored, err := repos.SchemaRepo.GetSchemaByUID(schema.UID)
	assert.NoError(t, err)
	assert.Equal(t, m3, stored.Mtime)

	// schema mtime ahead of the clock: the next mtime is still greater
	future := time.Now().UnixMilli() + 100000
	schema = testutils.CreateSchema(repos.SchemaRepo, t, user, &types.Schema{
		UID:   uuid.NewString(),
		Name:  testutils.CreateName(10),
		Mtime: future,
	})
	assert.Equal(t, future+1, save(32))
	stored, err = repos.SchemaRepo.GetSchemaByUID(schema.UID)
	assert.NoError(t, err)

	// UpdateSchema with a stale mtime doesn't change it
	stored.Mtime = 1
	stored.Description = "changed"
	assert.NoError(t, repos.SchemaRepo.UpdateSchema(stored))
	stored, err = repos.SchemaRepo.GetSchemaByUID(schema.UID)
	assert.NoError(t, err)
	assert.Equal(t, future+1, stored.Mtime)
	assert.Equal(t, "changed", stored.Description)

	// errors roll back the mtime and the changes
	fn_err := errors.New("fail")
	err = repos.SchemaPartRepo.UpdateWithNextMtime(schema.UID, func(repo *db.SchemaPartRepository, m int64) error {
		part := createPart(schema, 48)
		part.Mtime = m
		assert.NoError(t, repo.CreateOrUpdateSchemaPart(part))
		return fn_err
	})
	assert.ErrorIs(t, err, fn_err)
	stored, err = repos.SchemaRepo.GetSchemaByUID(schema.UID)
	assert.NoError(t, err)
	assert.Equal(t, future+1, stored.Mtime)
	part, err := repos.SchemaPartRepo.GetBySchemaUIDAndOffset(schema.UID, 48, 0, 0)
	assert.NoError(t, err)
	assert.Nil(t, part)

	// unknown schema
	err = repos.SchemaPartRepo.UpdateWithNextMtime("00000000-0000-0000-0000-000000000000", func(_ *db.SchemaPartRepository, _ int64) error {
		return nil
	})
	assert.Error(t, err)
}

func TestUpdateWithNextMtimeConcurrent(t *testing.T) {
	repos := CreateTestDatabase(t)
	user := testutils.CreateUser(repos.UserRepo, t, nil)
	schema := testutils.CreateSchema(repos.SchemaRepo, t, user, nil)
	start, err := repos.SchemaRepo.GetSchemaByUID(schema.UID)
	assert.NoError(t, err)

	// concurrent writes of different parts
	count := 20
	wg := sync.WaitGroup{}
	for i := 0; i < count; i++ {
		wg.Add(1)
		go func(x int) {
			defer wg.Done()
			err := repos.SchemaPartRepo.UpdateWithNextMtime(schema.UID, func(repo *db.SchemaPartRepository, m int64) error {
				part := createPart(schema, x*16)
				part.Mtime = m
				return repo.CreateOrUpdateSchemaPart(part)
			})
			assert.NoError(t, err)
		}(i)
	}
	wg.Wait()

	// paging by mtime returns every part exactly once
	seen := map[int]bool{}
	mtime := start.Mtime
	for {
		part, err := repos.SchemaPartRepo.GetNextBySchemaUIDAndMtime(schema.UID, mtime)
		assert.NoError(t, err)
		if part == nil {
			break
		}
		assert.Greater(t, part.Mtime, mtime)
		assert.False(t, seen[part.OffsetX], "part returned twice: %d", part.OffsetX)
		seen[part.OffsetX] = true
		mtime = part.Mtime
	}
	assert.Equal(t, count, len(seen))
}
