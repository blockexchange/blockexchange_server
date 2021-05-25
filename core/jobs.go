package core

import (
	"fmt"
	"os"
)

func StartJobs() {
	redis_host := os.Getenv("REDIS_HOST")
	redis_port := os.Getenv("REDIS_PORT")

	var syncmanager SyncManager
	if redis_host != "" && redis_port != "" {
		syncmanager = NewRedisSyncManager(redis_host + ":" + redis_port)
	} else {
		syncmanager = NewLocalSyncManager()
	}

	// TODO: add cleanup/screenshot-update jobs
	fmt.Println(syncmanager)
}
