package main

import (
	"blockexchange/api"
	"blockexchange/core"
	"blockexchange/db"
	"blockexchange/jobs"
	"blockexchange/types"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetLevel(logrus.InfoLevel)
	if os.Getenv("LOGLEVEL") == "debug" {
		logrus.SetLevel(logrus.DebugLevel)
	}

	logrus.Info("Starting")
	repos, err := db.Init()
	if err != nil {
		panic(err)
	}

	// populate database with test data (users, tokens)
	if os.Getenv("BLOCKEXCHANGE_TEST_DATA") == "true" {
		err = db.PopulateTestData(repos)
		if err != nil {
			panic(err)
		}
	}

	// set up server
	cfg := types.CreateConfig()

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort),
		Password: "",
		DB:       0,
	})

	api, router, err := api.NewApi(repos, cfg, rdb)
	if err != nil {
		panic(err)
	}

	// main entry
	http.Handle("/", router)
	server := &http.Server{Addr: ":8080", Handler: nil}

	// start background jobs
	jobs.Start(repos, api, core.New(cfg, repos))

	go func() {
		// listen to web requests
		err = server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	var captureSignal = make(chan os.Signal, 1)
	signal.Notify(captureSignal, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-captureSignal
	logrus.Info("Preparing shutdown")

	//stop api
	api.Stop()
	time.Sleep(5 * time.Second)
	logrus.Info("Shutdown complete")
	server.Shutdown(context.Background())
}
