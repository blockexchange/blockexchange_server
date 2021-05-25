package main

import (
	"blockexchange/db"
	"blockexchange/jobs"
	"blockexchange/web"
	"os"

	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetLevel(logrus.TraceLevel)
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

	// listen to web requests
	web.Serve(db_)
}
