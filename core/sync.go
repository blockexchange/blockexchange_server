package core

import (
	"sync"

	goredislib "github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
	"github.com/sirupsen/logrus"
)

type SyncManager interface {
	GetLock(name string) sync.Locker
}

// Local impl

func NewLocalSyncManager() SyncManager {
	return &LocalSyncManager{}
}

type LocalSyncManager struct {
}

func (sm *LocalSyncManager) GetLock(name string) sync.Locker {
	return &sync.Mutex{}
}

// Redis impl

func NewRedisSyncManager(redisUrl string) SyncManager {
	logrus.Infof("Syncmanager connecting to redis @ %s", redisUrl)
	client := goredislib.NewClient(&goredislib.Options{
		Addr: redisUrl,
	})

	pool := goredis.NewPool(client)
	rs := redsync.New(pool)

	return &RedisSyncManager{
		rs: rs,
	}
}

type RedisSyncManager struct {
	rs *redsync.Redsync
}

func (sm *RedisSyncManager) GetLock(name string) sync.Locker {
	return &RedisLockerWrapper{
		mutex: sm.rs.NewMutex(name),
	}
}

type RedisLockerWrapper struct {
	mutex *redsync.Mutex
}

func (rdw *RedisLockerWrapper) Lock() {
	err := rdw.mutex.Lock()
	if err != nil {
		panic(err)
	}
}

func (rdw *RedisLockerWrapper) Unlock() {
	_, err := rdw.mutex.Unlock()
	if err != nil {
		panic(err)
	}
}
