package main

import (
	"blockexchange/db"
	"blockexchange/web"
	"os"

	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetLevel(logrus.TraceLevel)
	logrus.SetReportCaller(true)
	logrus.Info("Starting")
	db_, err := db.Init()
	if err != nil {
		panic(err)
	}
	db.Migrate(db_.DB)

	if os.Getenv("BLOCKEXCHANGE_TEST_DATA") == "true" {
		// populate database with test data (users, tokens)
		err = db.PopulateTestData(db_)
		if err != nil {
			panic(err)
		}
	}
	web.Serve(db_)
}
