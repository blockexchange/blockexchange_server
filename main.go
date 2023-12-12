package main

import (
	"blockexchange/api"
	"blockexchange/db"
	"blockexchange/jobs"
	"blockexchange/types"
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetLevel(logrus.InfoLevel)
	logrus.Info("Starting")
	db_, err := db.Init()
	if err != nil {
		panic(err)
	}

	// migrate database
	db.Migrate(db_.DB)

	// populate database with test data (users, tokens)
	if os.Getenv("BLOCKEXCHANGE_TEST_DATA") == "true" {
		err = db.PopulateTestData(db_)
		if err != nil {
			panic(err)
		}
	}

	// start background jobs
	jobs.Start(db_)

	// set up server
	cfg := types.CreateConfig()
	api, err := api.NewApi(db_, cfg)
	if err != nil {
		panic(err)
	}
	server := &http.Server{Addr: ":8080", Handler: nil}

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
